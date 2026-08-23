package nts

import "net"

func NegotiationResponse(host string, port uint16) []Record {
	return []Record{{Type: RecordNextProtocol, Body: []byte{0, 0}}, {Type: RecordAEAD, Body: []byte{0, 15}}, {Type: RecordServer, Body: net.ParseIP(host).To16()}, {Type: RecordPort, Body: []byte{byte(port >> 8), byte(port)}}}
}
func IsCriticalUnsupported(r Record) bool {
	switch r.Type {
	case RecordNextProtocol, RecordAEAD, RecordServer, RecordPort, RecordEnd:
		return false
	default:
		return r.Critical
	}
}
