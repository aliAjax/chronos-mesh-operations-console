package source

import (
	"sync"
	"testing"
	"time"
)

func TestManagerListSnapshotIsolationC001(t *testing.T) {
	m := NewManager(nil)
	m.samples["a"] = Sample{Name: "a"}
	first := m.List()
	m.mu.Lock()
	delete(m.samples, "a")
	m.samples["b"] = Sample{Name: "b"}
	m.mu.Unlock()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; _ = m.List() }()
	}
	close(start)
	wg.Wait()
	if first[0].Name != "a" {
		t.Fatalf("first manager snapshot was mutated: %#v", first)
	}
}

func TestHistorySinceSnapshotIsolationC001(t *testing.T) {
	h := NewHistory(8)
	h.Add(Sample{Name: "old", Last: time.Unix(1, 0)})
	h.Add(Sample{Name: "new", Last: time.Unix(2, 0)})
	old := h.Since(time.Time{})
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; _ = h.Since(time.Unix(1, 0)) }()
	}
	close(start)
	wg.Wait()
	if old[0].Name != "old" {
		t.Fatalf("history snapshot was mutated: %#v", old)
	}
}

func TestEventLogSnapshotIsolationC001(t *testing.T) {
	log := NewEventLog(8)
	log.Append("primary", "up", "ready")
	log.Append("standby", "up", "ready")
	first := log.After(0)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; _ = log.After(1) }()
	}
	close(start)
	wg.Wait()
	if len(first) != 2 || first[0].Name != "primary" {
		t.Fatalf("event snapshot was mutated: %#v", first)
	}
}

func TestExplainConcurrentSnapshotsC001(t *testing.T) {
	start := make(chan struct{})
	results := make(chan string, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		results <- Explain(Sample{Name: "alpha", Address: "a", Accepted: true})
	}()
	go func() {
		defer wg.Done()
		<-start
		results <- Explain(Sample{Name: "beta", Address: "b", Accepted: false})
	}()
	close(start)
	wg.Wait()
	close(results)
	for result := range results {
		if result == "" {
			t.Fatal("empty diagnostic snapshot")
		}
	}
}
