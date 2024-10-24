package strategy

import (
	"loadBalancer/backend"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

type loadBalancer interface {
	nextIndex() int
	getCurrent() *uint64
	getBackends() []*backend.Backend
	setBackends([]*backend.Backend)
}

func HealthCheck(lb loadBalancer) {
	for _, b := range lb.getBackends() {
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
	}
}

func getNextPeer(lb loadBalancer) *backend.Backend {
	backends := lb.getBackends()
	next := lb.nextIndex()
	l := len(backends) + next

	for i := next; i < l; i++ {
		idx := i % len(backends)
		if backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(lb.getCurrent(), uint64(idx))
			}
			return backends[idx]
		}
	}
	return nil
}

func AddBackends(lb loadBalancer, b *backend.Backend) {
	backends := lb.getBackends()
	backends = append(backends, b)

	lb.setBackends(backends)
}

func GetLoadBalancer(balancer string) loadBalancer {
	switch balancer {
	case "RoundRobin":
		return &RoundRobin{}
	default:
		return nil
	}
}

func Serve(lb loadBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peer := getNextPeer(lb)
		if peer != nil {
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Service not available", http.StatusInternalServerError)
	}
}
