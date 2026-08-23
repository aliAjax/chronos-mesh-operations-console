package selection

import (
	"sync"
	"time"
)

type Hysteresis struct {
	mu      sync.Mutex
	last    time.Time
	minStay time.Duration
	offset  time.Duration
}

func NewHysteresis(d time.Duration) *Hysteresis { return &Hysteresis{minStay: d} }
func (h *Hysteresis) Accept(now time.Time, next time.Duration) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.last.IsZero() || now.Sub(h.last) >= h.minStay || absHysteresis(next-h.offset) > 100*time.Millisecond {
		h.offset = next
		h.last = now
		return true
	}
	return false
}
func absHysteresis(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
