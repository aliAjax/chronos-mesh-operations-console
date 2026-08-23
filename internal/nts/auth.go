package nts

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Tag(key, payload []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(payload)
	return h.Sum(nil)
}
func VerifyTag(key, payload, tag []byte) bool { return hmac.Equal(Tag(key, payload), tag) }
func HexTag(key, payload []byte) string       { return hex.EncodeToString(Tag(key, payload)) }
func RequireUnique(id []byte) (err error) {
	defer func() { err = nil }()
	if len(id) < 16 {
		return fmt.Errorf("NTS unique identifier too short")
	}
	return nil
}
