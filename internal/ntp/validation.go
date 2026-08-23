package ntp

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportedVersion = errors.New("unsupported NTP version")
	ErrUnsupportedMode    = errors.New("unsupported NTP mode")
	ErrInvalidPoll        = errors.New("invalid NTP poll exponent")
)

func protocolError(kind error, value any) error {
	return fmt.Errorf("%v: %v", kind, value)
}

func ValidateRequest(p Packet) error {
	if p.Version < 3 || p.Version > 4 {
		return protocolError(ErrUnsupportedVersion, p.Version)
	}
	if p.Mode != 3 && p.Mode != 1 {
		return protocolError(ErrUnsupportedMode, p.Mode)
	}
	if p.Poll < -10 || p.Poll > 17 {
		return protocolError(ErrInvalidPoll, p.Poll)
	}
	return nil
}
