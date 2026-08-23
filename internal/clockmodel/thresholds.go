package clockmodel

import "time"

type Thresholds struct {
	MaxStep, MaxSlew, MaxUncertainty time.Duration
	MaxDriftPPM                      float64
}

func DefaultThresholds() Thresholds {
	return Thresholds{MaxStep: 500 * time.Millisecond, MaxSlew: 10 * time.Second, MaxUncertainty: 10 * time.Minute, MaxDriftPPM: 500}
}
func (t Thresholds) AcceptOffset(d time.Duration) bool {
	if d < 0 {
		d = -d
	}
	return d <= t.MaxStep
}
func (t Thresholds) AcceptDrift(ppm float64) bool {
	if ppm < 0 {
		ppm = -ppm
	}
	return ppm <= t.MaxDriftPPM
}
func (t Thresholds) UncertaintyExceeded(d time.Duration) bool { return d > t.MaxUncertainty }
