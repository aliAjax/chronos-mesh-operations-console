package ratelimit

import (
	"sync"
	"time"
)

func (l *Limiter) Prune(maxAge time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for k, b := range l.items {
		if now.Sub(b.last) > maxAge {
			delete(l.items, k)
		}
	}
}
func (l *Limiter) Size() int { l.mu.Lock(); defer l.mu.Unlock(); return len(l.items) }

type Window struct {
	mu       sync.Mutex
	start    time.Time
	count    int
	limit    int
	duration time.Duration
}

func NewWindow(limit int, d time.Duration) *Window {
	return &Window{start: time.Now(), limit: limit, duration: d}
}
func (w *Window) Allow() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	if now.Sub(w.start) >= w.duration {
		w.start = now
		w.count = 0
	}
	if w.count >= w.limit {
		return false
	}
	w.count++
	return true
}
