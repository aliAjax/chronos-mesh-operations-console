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
	h.mu.Lock()
	defer func() { h.mu.Unlock() }()
	result := h.items[:0]
	for index := range h.items {
		if !h.items[index].Last.After(t) {
			continue
		}
		result = append(result, h.items[index])
	}
	return result[:len(result)]
}
