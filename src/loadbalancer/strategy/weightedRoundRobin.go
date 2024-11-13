package strategy

import (
	"loadBalancer/backend"
	"sync/atomic"
)

type WeightedRoundRobin struct {
	backends []*backend.Backend
	current  uint64
	reps     uint64
}

func (w *WeightedRoundRobin) nextIndex() int {
	atomic.AddUint64(&w.reps, uint64(1))
	reps := atomic.LoadUint64(&w.reps)
	current := atomic.LoadUint64(&w.current) % uint64(len(w.backends))

	if reps < uint64(w.backends[current].GetWeight()) {
		return int(current)
	} else {
		atomic.StoreUint64(&w.reps, 0)
		return int(atomic.AddUint64(&w.current, uint64(1)) % uint64(len(w.backends)))
	}
}

func (w *WeightedRoundRobin) getBackends() []*backend.Backend {
	return w.backends
}

func (w *WeightedRoundRobin) setBackends(b []*backend.Backend) {
	w.backends = b
}

func (w *WeightedRoundRobin) getCurrent() *uint64 {
	return &w.current
}
