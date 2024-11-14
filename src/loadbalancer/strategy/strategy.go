package strategy

import (
	"encoding/json"
	"loadBalancer/backend"
	"log"
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

func MarkBackendStateus(u *url.URL, status bool, lb loadBalancer) {
	for _, b := range lb.getBackends() {
		if b.URL.String() == u.String() {
			b.SetAlive(status)
			break
		}
	}
}

func BackendStatusHandler(lb loadBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := []*struct {
			URL   string `json:"url"`
			Alive bool `json:"alive"`
		}{}
		for _, b := range lb.getBackends() {
			response = append(response, &struct {
        URL   string `json:"url"`
        Alive bool `json:"alive"`
			}{
				URL:   b.URL.String(),
				Alive: b.IsAlive(),
			})
		}

    responseJson, err := json.Marshal(response)
    if err != nil {
      http.Error(w, "Error creating health json", http.StatusInternalServerError)
    }

    w.Write(responseJson)
	}
}

func HealthCheck(lb loadBalancer) {
	for _, b := range lb.getBackends() {
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)

    log.Printf("Backend %s is alive %t", b.URL.String(), alive)
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
  case "WeightedRoundRobin":
    return &WeightedRoundRobin{}
  case "Random":
    return &Random{}
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
