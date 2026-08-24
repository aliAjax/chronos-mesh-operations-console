package ratelimit

import (
	"net"
	"sync"
	"time"
)

type bucket struct {
	tokens float64
	last   time.Time
}
type Limiter struct {
	mu    sync.Mutex
	rate  float64
	burst float64
	items map[string]bucket
	Stats
}

func New(rate int) *Limiter {
	return &Limiter{rate: float64(rate), burst: float64(rate), items: map[string]bucket{}}
}
func (l *Limiter) Allow(ip net.IP) bool {
	key := ip.String()
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	b := l.items[key]
	if b.last.IsZero() {
		b.last = now
		b.tokens = l.burst
	}
	b.tokens += now.Sub(b.last).Seconds() * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now
	if b.tokens < 1 {
		l.Denied.Add(1)
		return false
	}
	b.tokens--
	l.items[key] = b
	l.Allowed.Add(1)
	return true
}
