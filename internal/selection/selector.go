package selection

import (
	"github.com/chronos-mesh/chronos-mesh/internal/source"
	"sort"
	"time"
)

type Decision struct {
	Offset  time.Duration
	Sources []source.Sample
	Trusted int
	Reason  string
}

func Select(in []source.Sample) Decision {
	c := append([]source.Sample(nil), in...)
	sort.Slice(c, func(i, j int) bool { return c[i].Offset < c[j].Offset })
	valid := make([]source.Sample, 0, len(c))
	for _, s := range c {
		if s.Reach == 0 || s.Stratum == 0 || s.Delay < 0 || s.Delay > 2*time.Second {
			s.Accepted = false
			s.Reason = "invalid-delay-or-stratum"
			continue
		}
		valid = append(valid, s)
	}
	if len(valid) == 0 {
		return Decision{Sources: c, Reason: "no-trusted-source"}
	}
	med := valid[len(valid)/2].Offset
	kept := valid[:0]
	for _, s := range valid {
		if abs(s.Offset-med) <= 250*time.Millisecond {
			s.Accepted = true
			s.Reason = "inlier"
			kept = append(kept, s)
		} else {
			s.Accepted = false
			s.Reason = "falseticker"
		}
	}
	var sum time.Duration
	for _, s := range kept {
		sum += s.Offset
	}
	if len(kept) > 0 {
		sum /= time.Duration(len(kept))
	}
	return Decision{Offset: sum, Sources: c, Trusted: len(kept), Reason: "clustered"}
}
func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
