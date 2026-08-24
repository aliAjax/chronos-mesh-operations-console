package nts

import (
	"sync"
	"testing"
	"time"
)

func TestReplayAcceptPublishesAtomically(t *testing.T) {
	r := NewReplay(32, time.Minute)
	var wg sync.WaitGroup
	var mu sync.Mutex
	wins := 0
	start := make(chan struct{})
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if r.Accept("same", time.Now()) {
				mu.Lock()
				wins++
				mu.Unlock()
			}
		}()
	}
	close(start)
	wg.Wait()
	if wins != 1 {
		t.Fatalf("replay accepted %d concurrent copies", wins)
	}
}

func TestRotatedKeyCopyStaysDetached(t *testing.T) {
	k := NewKeyRing([]byte("seed"), time.Minute)
	snapshot := k.Current()
	snapshot[0] ^= 0xff
	if snapshot[0] == k.Current()[0] {
		t.Fatal("mutating key snapshot changed key ring")
	}
}

func TestKEServerHandlerSnapshotConcurrent(t *testing.T) {
	s := &KEServer{}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); s.SetHandler(nil) }()
		go func() { defer wg.Done(); _ = s.Handler() }()
	}
	wg.Wait()
}

func TestNegotiatedCapabilitiesStayDetached(t *testing.T) {
	original := DefaultCapabilities()
	snapshot := original.Snapshot()
	snapshot.Protocols[0] = 99
	snapshot.AEADs[0] = 98
	if original.Protocols[0] == 99 || original.AEADs[0] == 98 {
		t.Fatal("capability snapshot aliases internal slices")
	}
}
