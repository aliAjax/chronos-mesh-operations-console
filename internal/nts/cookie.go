package nts

import (
	"crypto/rand"
	"fmt"
	"time"
)

type Cookie struct {
	ID      []byte
	Expires time.Time
}

func NewID() ([]byte, error) { b := make([]byte, 32); _, e := rand.Read(b); return b, e }
func EncodeCookie(m *Manager, id []byte) ([]byte, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("empty cookie id")
	}
	return m.Cookie(id)
}
func DecodeCookie(m *Manager, b []byte) (Cookie, error) {
	id, e := m.Open(b)
	if e != nil {
		return Cookie{}, e
	}
	return Cookie{ID: id, Expires: time.Now().Add(m.ttl)}, nil
}
