package ntp

import (
	"fmt"
	"time"
)

type Diagnostic struct {
	Version  uint8
	Mode     uint8
	Stratum  uint8
	Leap     string
	Transmit time.Time
	Valid    bool
	Error    string
}

func Inspect(b []byte) Diagnostic {
	p, e := Decode(b)
	if e != nil {
		return Diagnostic{Valid: false, Error: e.Error()}
	}
	e = ValidateRequest(p)
	return Diagnostic{Version: p.Version, Mode: p.Mode, Stratum: p.Stratum, Leap: LeapIndicator(p.LI), Transmit: p.Transmit, Valid: e == nil, Error: errorText(e)}
}
func errorText(e error) string {
	if e == nil {
		return ""
	}
	return e.Error()
}
func Summary(p Packet) string {
	return fmt.Sprintf("v=%d mode=%d stratum=%d leap=%s transmit=%s", p.Version, p.Mode, p.Stratum, LeapIndicator(p.LI), p.Transmit.Format(time.RFC3339Nano))
}
func IsClientRequest(p Packet) bool  { return p.Mode == 3 || p.Mode == 1 }
func IsServerResponse(p Packet) bool { return p.Mode == 4 || p.Mode == 2 }
