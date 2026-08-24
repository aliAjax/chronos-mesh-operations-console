package nts

import (
	"crypto/sha256"
	"sync"
	"time"
)

type KeyRing struct {
	mu       sync.RWMutex
	current  []byte
	previous []byte
	changed  time.Time
	period   time.Duration
}

type KeySnapshot struct {
	Current  []byte
	Previous []byte
	Changed  time.Time
}

func NewKeyRing(seed []byte, period time.Duration) *KeyRing {
	h := sha256.Sum256(seed)
	return &KeyRing{current: cloneBytes(h[:]), changed: time.Now(), period: period}
}

// Rotate advances the key ring: the current key becomes the previous one and a
// new current key is derived from seed. Both keys are independent copies, so a
// caller that still holds a prior Current/Snapshot value cannot observe this
// rotation mutating its bytes.
func (k *KeyRing) Rotate(seed []byte) {
	h := sha256.Sum256(seed)
	k.mu.Lock()
	k.previous = k.current
	k.current = cloneBytes(h[:])
	k.changed = time.Now()
	k.mu.Unlock()
}

// Current returns a defensive copy of the active key. The returned slice is
// independent of the ring and of any other caller's snapshot, so mutating it
// cannot corrupt another request's view of the key.
func (k *KeyRing) Current() []byte {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return cloneBytes(k.current)
}

// Snapshot atomically returns independent copies of the current and previous
// keys plus the rotation time. Use this when you need a consistent pair (e.g.
// to verify a nonce against either the current or the just-rotated previous
// key without a race against concurrent Rotate calls).
func (k *KeyRing) Snapshot() KeySnapshot {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return KeySnapshot{
		Current:  cloneBytes(k.current),
		Previous: cloneBytes(k.previous),
		Changed:  k.changed,
	}
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// Snapshot returns an independent copy of the capability set so that later
// renegotiation (which may replace the Protocols/AEADs slices in place) cannot
// mutate a snapshot already handed to a caller. The slice backing arrays are
// duplicated, not shared.
func (c CapabilitySet) Snapshot() CapabilitySet {
	return CapabilitySet{
		Protocols:    cloneUint16s(c.Protocols),
		AEADs:        cloneUint16s(c.AEADs),
		CookieLength: c.CookieLength,
		ReplayWindow: c.ReplayWindow,
	}
}

func cloneUint16s(s []uint16) []uint16 {
	if s == nil {
		return nil
	}
	c := make([]uint16, len(s))
	copy(c, s)
	return c
}
