package strategy

import (
	"loadBalancer/backend"
	"math/rand"
)

type Random struct {
	backends []*backend.Backend
	current  uint64
}

func (r *Random) nextIndex() int {
	return rand.Intn(3)
}

func (r *Random) getBackends() []*backend.Backend {
  return r.backends
}

func (r *Random) setBackends(b []*backend.Backend) {
  r.backends = b
}

func (r *Random) getCurrent() *uint64 {
  return &r.current
}
