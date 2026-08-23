package source

import (
	"fmt"
	"strings"
	"time"
)

func Explain(s Sample) string {
	parts := []string{s.Name, s.Address}
	if s.Accepted {
		parts = append(parts, "accepted")
	} else {
		parts = append(parts, "rejected")
	}
	if s.Reason != "" {
		parts = append(parts, s.Reason)
	}
	parts = append(parts, fmt.Sprintf("offset=%s", s.Offset), fmt.Sprintf("delay=%s", s.Delay))
	return strings.Join(parts, " ")
}
func Age(s Sample, now time.Time) time.Duration {
	if s.Last.IsZero() {
		return time.Duration(1<<63 - 1)
	}
	return now.Sub(s.Last)
}
func Healthy(s Sample, now time.Time, maxAge time.Duration) bool {
	return s.Accepted && Age(s, now) <= maxAge
}
