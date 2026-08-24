package ntp

import "time"

type FixedPoint uint32

func EncodeDuration(d time.Duration) FixedPoint {
	if d < 0 {
		return 0
	}
	v := float64(d) / float64(time.Second) * 65536
	if v > 4294967295 {
		return ^FixedPoint(0)
	}
	return FixedPoint(v)
}
func DecodeDuration(v FixedPoint) time.Duration {
	return time.Duration(float64(v) / 65536 * float64(time.Second))
}
func LeapIndicator(li uint8) string {
	switch li {
	case 0:
		return "no-warning"
	case 1:
		return "last-minute-61"
	case 2:
		return "last-minute-59"
	default:
		return "unsynchronized"
	}
}
