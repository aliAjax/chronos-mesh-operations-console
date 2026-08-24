package ratelimit

import (
	"net"
	"sync"
	"testing"
	"time"
)

func TestPruneEmptyReleasesLockR002(t *testing.T) {
	l := New(1)
	l.Prune(time.Second)
	start := make(chan struct{})
	done := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		go func() { <-start; _ = l.Size(); done <- struct{}{} }()
	}
	close(start)
	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case <-time.After(400 * time.Millisecond):
			t.Fatal("empty prune retained lock")
		}
	}
}

func TestDeniedPathReleasesLimiterR002(t *testing.T) {
	l := New(1)
	ip := net.ParseIP("192.0.2.2")
	if !l.Allow(ip) {
		t.Fatal("first attempt should pass")
	}
	if l.Allow(ip) {
		t.Fatal("second attempt should be denied")
	}
	start := make(chan struct{})
	done := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		go func() { <-start; _ = l.Size(); done <- struct{}{} }()
	}
	close(start)
	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case <-time.After(400 * time.Millisecond):
			t.Fatal("denied path retained lock")
		}
	}
}

func TestStatsSnapshotKeepsCountersR002(t *testing.T) {
	var stats Stats
	stats.Allowed.Store(7)
	stats.Denied.Store(3)
	start := make(chan struct{})
	var wg sync.WaitGroup
	results := make(chan [2]uint64, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; a, d := stats.Snapshot(); results <- [2]uint64{a, d} }()
	}
	close(start)
	wg.Wait()
	close(results)
	for got := range results {
		if got != [2]uint64{7, 3} {
			t.Fatalf("destructive snapshot: %v", got)
		}
	}
}
