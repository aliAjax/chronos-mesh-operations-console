package nts

import (
	"sync"
	"time"
)

type Replay struct {
	mu   sync.Mutex
	seen map[string]time.Time
	ttl  time.Duration
	max  int
}

func NewReplay(max int, ttl time.Duration) *Replay {
	return &Replay{seen: map[string]time.Time{}, max: max, ttl: ttl}
}
func (r *Replay) Accept(id string, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, t := range r.seen {
		if now.Sub(t) > r.ttl {
			delete(r.seen, k)
		}
	}
	if _, ok := r.seen[id]; ok {
		return false
	}
	if len(r.seen) >= r.max {
		for k := range r.seen {
			delete(r.seen, k)
			break
		}
	}
	r.seen[id] = now
	return true
}
