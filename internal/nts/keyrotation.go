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

func NewKeyRing(seed []byte, period time.Duration) *KeyRing {
	h := sha256.Sum256(seed)
	return &KeyRing{current: h[:], changed: time.Now(), period: period}
}
func (k *KeyRing) Rotate(seed []byte) {
	h := sha256.Sum256(seed)
	k.mu.Lock()
	k.previous = k.current
	k.current = h[:]
	k.changed = time.Now()
	k.mu.Unlock()
}
func (k *KeyRing) Current() []byte {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return append([]byte(nil), k.current...)
}
