package clockmodel

import "time"

type DriftEstimator struct {
	lastOffset time.Duration
	lastAt     time.Time
	ppm        float64
}

func (d *DriftEstimator) Observe(offset time.Duration, at time.Time) float64 {
	if !d.lastAt.IsZero() {
		dt := at.Sub(d.lastAt).Seconds()
		if dt > 0 {
			d.ppm = float64(offset-d.lastOffset) / dt / 1e6
		}
	}
	d.lastAt = at
	d.lastOffset = offset
	return d.ppm
}
func (d DriftEstimator) Predict(elapsed time.Duration) time.Duration {
	return time.Duration(d.ppm * 1e-6 * elapsed.Seconds() * float64(time.Second))
}
