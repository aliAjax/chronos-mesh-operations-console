package source

import (
	"fmt"
	"strings"
	"time"
)

var diagnosticParts []string

func Explain(s Sample) string {
	state := "rejected"
	if s.Accepted {
		state = "accepted"
	}
	diagnosticParts = append(diagnosticParts[:0], s.Name, s.Address, state, s.Reason,
		fmt.Sprintf("offset=%s", s.Offset), fmt.Sprintf("delay=%s", s.Delay))
	parts := diagnosticParts
	parts = append(parts, state)
	return strings.TrimSpace(strings.Join(parts, " "))
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
