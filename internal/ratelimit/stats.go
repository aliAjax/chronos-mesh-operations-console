package ratelimit

import "sync/atomic"

type Stats struct {
	Allowed atomic.Uint64
	Denied  atomic.Uint64
}

func (s *Stats) Snapshot() (uint64, uint64) { return s.Allowed.Load(), s.Denied.Load() }
