package ntp

import "fmt"

func ValidateRequest(p Packet) error {
	if p.Version < 3 || p.Version > 4 {
		return fmt.Errorf("unsupported NTP version %d", p.Version)
	}
	if p.Mode != 3 && p.Mode != 1 {
		return fmt.Errorf("unsupported mode %d", p.Mode)
	}
	if p.Poll < -10 || p.Poll > 17 {
		return fmt.Errorf("invalid poll exponent")
	}
	return nil
}
