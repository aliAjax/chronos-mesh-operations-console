package control

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

func APIKey(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expected == "" {
				next.ServeHTTP(w, r)
				return
			}
			got := r.Header.Get("Authorization")
			a := sha256.Sum256([]byte(got))
			b := sha256.Sum256([]byte(expected))
			if subtle.ConstantTimeCompare(a[:], b[:]) != 1 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
