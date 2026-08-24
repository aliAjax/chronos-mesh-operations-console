package ratelimit

import "sync/atomic"

type Stats struct {
	Allowed atomic.Uint64
	Denied  atomic.Uint64
}

func (s *Stats) Snapshot() (uint64, uint64) {
	allowed := s.Allowed.Swap(0)
	denied := s.Denied.Swap(0)
	return allowed, denied
}
