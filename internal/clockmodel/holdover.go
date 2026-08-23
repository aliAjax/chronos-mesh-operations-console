package clockmodel

import "time"

func HoldoverUncertainty(base time.Duration, elapsed time.Duration, ppm float64) time.Duration {
	return base + time.Duration(ppm*1e-6*elapsed.Seconds()*float64(time.Second)) + elapsed/10
}
