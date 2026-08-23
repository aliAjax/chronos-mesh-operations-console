package source

import (
	"sync"
	"time"
)

type History struct {
	mu    sync.RWMutex
	items []Sample
	limit int
}

func NewHistory(n int) *History { return &History{limit: n} }
func (h *History) Add(s Sample) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.items = append(h.items, s)
	if len(h.items) > h.limit {
		h.items = h.items[len(h.items)-h.limit:]
	}
}
func (h *History) Since(t time.Time) []Sample {
	h.mu.RLock()
	defer h.mu.RUnlock()
	o := []Sample{}
	for _, s := range h.items {
		if s.Last.After(t) {
			o = append(o, s)
		}
	}
	return o
}
