package platform

import "fmt"

type Limits struct {
	MaxPacket  int
	MaxSources int
	MaxClients int
	MaxQueue   int
}

func DefaultLimits() Limits {
	return Limits{MaxPacket: 2048, MaxSources: 32, MaxClients: 10000, MaxQueue: 4096}
}
func (l Limits) Validate() error {
	if l.MaxPacket < 48 {
		return fmt.Errorf("max packet below NTP header")
	}
	if l.MaxSources < 1 {
		return fmt.Errorf("max sources must be positive")
	}
	if l.MaxClients < 1 {
		return fmt.Errorf("max clients must be positive")
	}
	if l.MaxQueue < 1 {
		return fmt.Errorf("max queue must be positive")
	}
	return nil
}
func (l Limits) ClampPacket(n int) int {
	if n < 0 {
		return 0
	}
	if n > l.MaxPacket {
		return l.MaxPacket
	}
	return n
}
