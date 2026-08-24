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

// Accept atomically checks whether id has been seen and, if not, records it.
// The check and the insert run under a single critical section so that two
// concurrent calls with the same id cannot both pass (which would let a
// duplicated nonce through).
func (r *Replay) Accept(id string, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Evict expired entries before deciding so the window reflects reality.
	for k, t := range r.seen {
		if now.Sub(t) > r.ttl {
			delete(r.seen, k)
		}
	}

	if _, ok := r.seen[id]; ok {
		return false
	}

	// Bound the table: drop one expired-or-oldest entry when full.
	if len(r.seen) >= r.max {
		var oldestK string
		var oldestT time.Time
		first := true
		for k, t := range r.seen {
			if first || t.Before(oldestT) {
				oldestK, oldestT, first = k, t, false
			}
		}
		delete(r.seen, oldestK)
	}

	r.seen[id] = now
	return true
}
