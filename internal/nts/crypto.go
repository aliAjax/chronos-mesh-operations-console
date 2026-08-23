package nts

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

var ErrMalformedCookie = errors.New("malformed cookie")

type Manager struct {
	key  []byte
	aead cipher.AEAD
	ttl  time.Duration
}

func New(key string) (*Manager, error) {
	h := sha256.Sum256([]byte(key))
	b, e := aes.NewCipher(h[:])
	if e != nil {
		return nil, e
	}
	a, e := cipher.NewGCM(b)
	if e != nil {
		return nil, e
	}
	return &Manager{key: h[:], aead: a, ttl: time.Hour}, nil
}
func (m *Manager) Cookie(id []byte) ([]byte, error) {
	nonce := make([]byte, m.aead.NonceSize())
	if _, e := rand.Read(nonce); e != nil {
		return nil, e
	}
	plain := make([]byte, 8+len(id))
	binary.BigEndian.PutUint64(plain, uint64(time.Now().Add(m.ttl).Unix()))
	copy(plain[8:], id)
	return m.aead.Seal(nonce, nonce, plain, nil), nil
}
func (m *Manager) Open(cookie []byte) ([]byte, error) {
	n := m.aead.NonceSize()
	if len(cookie) < n {
		return nil, fmt.Errorf("cookie too short")
	}
	p, e := m.aead.Open(nil, cookie[:n], cookie[n:], nil)
	if e != nil {
		return nil, e
	}
	if int64(binary.BigEndian.Uint64(p)) < time.Now().Unix() {
		return nil, fmt.Errorf("cookie expired")
	}
	return p[8:], nil
}
func (m *Manager) Authenticate(payload, cookie, tag []byte) bool {
	h := sha256.New()
	h.Write(payload)
	h.Write(cookie)
	sum := h.Sum(nil)
	return len(tag) == len(sum) && string(tag) == string(sum)
}
