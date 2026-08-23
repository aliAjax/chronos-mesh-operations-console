package clockmodel

import (
	"github.com/chronos-mesh/chronos-mesh/internal/ntp"
	"github.com/chronos-mesh/chronos-mesh/internal/platform"
	"github.com/chronos-mesh/chronos-mesh/internal/selection"
	"github.com/chronos-mesh/chronos-mesh/internal/source"
	"sync/atomic"
	"time"
)

type Snapshot struct {
	Base                       time.Time
	Offset, Drift, Uncertainty time.Duration
	Stratum                    uint8
	Reference                  time.Time
	Synchronized               bool
	Holdover                   time.Duration
}
type Model struct {
	clock   platform.Clock
	maxHold time.Duration
	state   atomic.Value
}

func New(c platform.Clock, maxHold time.Duration) *Model {
	m := &Model{clock: c, maxHold: maxHold}
	m.state.Store(Snapshot{Base: c.Now(), Stratum: 16, Uncertainty: time.Hour})
	return m
}
func (m *Model) Update(d selection.Decision, samples []source.Sample) {
	now := m.clock.Now()
	old := m.state.Load().(Snapshot)
	synced := d.Trusted > 0
	str := uint8(16)
	if synced {
		str = 2
	}
	hold := time.Duration(0)
	if !synced {
		hold = now.Sub(old.Base)
		if hold > m.maxHold {
			str = 16
		}
	}
	u := old.Uncertainty
	if synced {
		u = 250 * time.Millisecond
	} else {
		u += time.Second
	}
	m.state.Store(Snapshot{Base: now, Offset: d.Offset, Drift: old.Drift, Uncertainty: u, Stratum: str, Reference: now, Synchronized: synced, Holdover: hold})
}
func (m *Model) Snapshot() ntp.Snapshot {
	s := m.state.Load().(Snapshot)
	now := m.clock.Now()
	elapsed := now.Sub(s.Base)
	return ntp.Snapshot{Now: now.Add(s.Offset), Stratum: s.Stratum, Reference: s.Reference, Synchronized: s.Synchronized, RootDelay: elapsed / 10, RootDispersion: s.Uncertainty + elapsed/100}
}
func (m *Model) State() Snapshot { return m.state.Load().(Snapshot) }
