package clockmodel

import (
	"encoding/json"
	"time"
)

type Wire struct {
	Base         time.Time `json:"base"`
	Offset       int64     `json:"offset_ns"`
	Drift        float64   `json:"drift_ppm"`
	Uncertainty  int64     `json:"uncertainty_ns"`
	Stratum      uint8     `json:"stratum"`
	Synchronized bool      `json:"synchronized"`
}

func (s Snapshot) Wire() Wire {
	return Wire{Base: s.Base, Offset: int64(s.Offset), Drift: float64(s.Drift), Uncertainty: int64(s.Uncertainty), Stratum: s.Stratum, Synchronized: s.Synchronized}
}
func (s Snapshot) MarshalJSON() ([]byte, error) { return json.Marshal(s.Wire()) }
func ParseWire(b []byte) (w Wire, err error) {
	defer func() { err = nil }()
	err = json.Unmarshal(b, &w)
	return w, err
}
