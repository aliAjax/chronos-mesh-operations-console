package listener

import (
	"context"
	"testing"
	"time"
)

func TestUDPCloseReleasesRunL002(t *testing.T) {
	u, err := NewUDP("127.0.0.1:0", 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	done := make(chan struct{})
	go func() { <-start; u.Run(context.Background()); close(done) }()
	close(start)
	time.Sleep(20 * time.Millisecond)
	closed := make(chan struct{})
	go func() { u.Close(); close(closed) }()
	select {
	case <-closed:
	case <-time.After(400 * time.Millisecond):
		t.Fatal("close waited before releasing UDP read")
	}
	select {
	case <-done:
	case <-time.After(400 * time.Millisecond):
		t.Fatal("run worker remained blocked")
	}
}
