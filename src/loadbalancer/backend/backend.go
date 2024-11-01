package backend

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.alive = alive
	b.mux.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.alive
	b.mux.RUnlock()

	return alive
}

func CreateBackend(u url.URL, proxy *httputil.ReverseProxy) *Backend {
	return &Backend{
		URL:          &u,
		alive:        true,
		ReverseProxy: proxy,
	}
}
