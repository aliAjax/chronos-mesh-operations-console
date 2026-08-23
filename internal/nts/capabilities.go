package nts

type CapabilitySet struct {
	Protocols    []uint16
	AEADs        []uint16
	CookieLength int
	ReplayWindow int
}

func DefaultCapabilities() CapabilitySet {
	return CapabilitySet{Protocols: []uint16{0}, AEADs: []uint16{15}, CookieLength: 64, ReplayWindow: 4096}
}
func (c CapabilitySet) SupportsProtocol(v uint16) bool {
	for _, x := range c.Protocols {
		if x == v {
			return true
		}
	}
	return false
}
func (c CapabilitySet) SupportsAEAD(v uint16) bool {
	for _, x := range c.AEADs {
		if x == v {
			return true
		}
	}
	return false
}
func (c CapabilitySet) Valid() bool {
	return len(c.Protocols) > 0 && len(c.AEADs) > 0 && c.CookieLength >= 32 && c.ReplayWindow > 0
}
