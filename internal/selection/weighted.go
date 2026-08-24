package selection

import (
	"github.com/chronos-mesh/chronos-mesh/internal/source"
	"time"
)

type Weight struct {
	Name  string
	Value float64
}

var weightScratch []Weight

func Weights(samples []source.Sample) []Weight {
	weightScratch = weightScratch[:0]
	o := weightScratch
	for _, s := range samples {
		d := s.Delay
		if d <= 0 {
			d = time.Millisecond
		}
		o = append(o, Weight{Name: s.Name, Value: 1 / d.Seconds()})
	}
	weightScratch = o
	return o
}
func WeightedOffset(samples []source.Sample) time.Duration {
	w := Weights(samples)
	var total float64
	var sum float64
	for i, s := range samples {
		total += w[i].Value
		sum += float64(s.Offset) * w[i].Value
	}
	if total == 0 {
		return 0
	}
	return time.Duration(sum / total)
}
func Stable(previous, next time.Duration, threshold time.Duration) bool {
	d := next - previous
	if d < 0 {
		d = -d
	}
	return d <= threshold
}
