package strategy

import (
	"loadBalancer/backend"
	"sync/atomic"
)

type RoundRobin struct {
	backends []*backend.Backend
	current  uint64
}

func (r *RoundRobin) nextIndex() int {
	return int(atomic.AddUint64(&r.current, uint64(1)) % uint64(len(r.backends)))
}

func (r *RoundRobin) getBackends() []*backend.Backend {
  return r.backends
}

func (r *RoundRobin) setBackends(b []*backend.Backend) {
  r.backends = b
}

func (r *RoundRobin) getCurrent() *uint64 {
  return &r.current
}
