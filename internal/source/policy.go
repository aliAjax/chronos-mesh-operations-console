package source

import "time"

type Policy struct {
	MaxDelay, MaxRootDistance, MaxOffset time.Duration
	MinReach                             uint8
	Hold                                 time.Duration
}

func DefaultPolicy() Policy {
	return Policy{MaxDelay: 2 * time.Second, MaxRootDistance: 5 * time.Second, MaxOffset: 2 * time.Second, MinReach: 1, Hold: time.Hour}
}
func (p Policy) Accept(s Sample) bool {
	return s.Reach >= p.MinReach && s.Delay >= 0 && s.Delay <= p.MaxDelay && s.RootDistance <= p.MaxRootDistance && abs(s.Offset) <= p.MaxOffset
}
func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
