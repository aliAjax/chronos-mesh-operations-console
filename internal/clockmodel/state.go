package clockmodel

import (
	"sync"
	"time"
)

type HistoryEntry struct {
	At     time.Time
	State  string
	Reason string
	Offset time.Duration
}
type History struct {
	mu    sync.RWMutex
	items []HistoryEntry
	limit int
}

func NewHistory(limit int) *History { return &History{limit: limit} }
func (h *History) Add(e HistoryEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.items = append(h.items, e)
	if len(h.items) > h.limit {
		h.items = h.items[len(h.items)-h.limit:]
	}
}
func (h *History) List() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return append([]HistoryEntry(nil), h.items...)
}
func (h *History) Last() HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.items) == 0 {
		return HistoryEntry{}
	}
	return h.items[len(h.items)-1]
}
