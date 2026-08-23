package clockmodel

import "time"

type Slew struct{ Rate float64 }

func (s Slew) Adjust(offset time.Duration, elapsed time.Duration) time.Duration {
	if elapsed <= 0 {
		return offset
	}
	max := time.Duration(float64(elapsed) * s.Rate)
	if offset > max {
		return max
	}
	if offset < -max {
		return -max
	}
	return offset
}
