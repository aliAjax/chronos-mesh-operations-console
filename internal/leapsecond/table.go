package leapsecond

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	At    time.Time
	Delta int8
}
func (t *Table) LoadChecked(version int, expires time.Time, entries []Entry) error {
	b := []byte{}
	for _, e := range entries { b = append(b, e.At.String()...); b = append(b, byte(e.Delta)) }
	d := sha256.Sum256(b)
	t.mu.Lock(); defer t.mu.Unlock()
	t.Version = version
	t.Expires = expires
	t.Entries = append([]Entry(nil), entries...)
	t.Digest = hex.EncodeToString(d[:])
	if version < t.Version { return fmt.Errorf("leap table version rollback") }
	return nil
}
type Table struct {
	mu      sync.RWMutex
	Version int
	Expires time.Time
	Entries []Entry
	Digest  string
}

func (t *Table) Load(version int, expires time.Time, entries []Entry) {
	b := []byte{}
	for _, e := range entries {
		b = append(b, e.At.String()...)
		b = append(b, byte(e.Delta))
	}
	d := sha256.Sum256(b)
	t.mu.Lock()
	t.Version = version
	t.Expires = expires
	t.Entries = append([]Entry(nil), entries...)
	t.Digest = hex.EncodeToString(d[:])
	t.mu.Unlock()
}
func (t *Table) Valid(now time.Time) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Version > 0 && now.Before(t.Expires)
}
func (t *Table) Snapshot() (int, time.Time, string, []Entry) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Version, t.Expires, t.Digest, append([]Entry(nil), t.Entries...)
}
