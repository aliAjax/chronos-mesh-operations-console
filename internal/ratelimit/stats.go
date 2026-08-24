package ratelimit

import "sync/atomic"

type Stats struct {
	Allowed atomic.Uint64
	Denied  atomic.Uint64
}

// Snapshot returns the cumulative allowed and denied counters without
// resetting them. It is safe for concurrent callers (e.g. /metrics and the
// operations panel polling together); the previous Swap(0) implementation
// destroyed counts on read, so a poll could observe swapped/zero values.
func (s *Stats) Snapshot() (uint64, uint64) {
	return s.Allowed.Load(), s.Denied.Load()
}
