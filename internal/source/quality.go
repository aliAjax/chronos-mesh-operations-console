package source

import (
	"math"
	"sort"
	"time"
)

func MedianOffset(samples []Sample) time.Duration {
	if len(samples) == 0 {
		return 0
	}
	x := append([]Sample(nil), samples...)
	sort.Slice(x, func(i, j int) bool { return x[i].Offset < x[j].Offset })
	return x[len(x)/2].Offset
}
func Jitter(samples []Sample, center time.Duration) time.Duration {
	if len(samples) == 0 {
		return 0
	}
	var sum float64
	for _, s := range samples {
		d := float64(s.Offset - center)
		sum += d * d
	}
	return time.Duration(math.Sqrt(sum / float64(len(samples))))
}
func RootDistance(s Sample) time.Duration { return s.Delay/2 + s.Jitter }
func ReachShift(reach uint8, ok bool) uint8 {
	if ok {
		return reach<<1 | 1
	}
	return reach << 1
}
func PollInterval(stability time.Duration) time.Duration {
	if stability < time.Second {
		return time.Second
	}
	if stability > time.Minute {
		return time.Minute
	}
	return stability
}
