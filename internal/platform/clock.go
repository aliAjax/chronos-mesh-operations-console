package platform

import "time"

type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}
type RealClock struct{}

func (RealClock) Now() time.Time        { return time.Now().UTC() }
func (RealClock) Sleep(d time.Duration) { time.Sleep(d) }

type FixedClock struct{ T time.Time }

func (f *FixedClock) Now() time.Time        { return f.T }
func (f *FixedClock) Sleep(d time.Duration) { f.T = f.T.Add(d) }
