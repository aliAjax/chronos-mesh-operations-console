package ntp

import "encoding/binary"

type KissCode [4]byte

var (
	KoD  = KissCode{'R', 'A', 'T', 'E'}
	DENY = KissCode{'D', 'E', 'N', 'Y'}
	RSTR = KissCode{'R', 'S', 'T', 'R'}
)

func Kiss(code KissCode, version uint8) Packet {
	return Packet{Version: version, Mode: 4, Stratum: 0, ReferenceID: binary.BigEndian.Uint32(code[:])}
}
func CodeString(id uint32) string {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, id)
	return string(b)
}
func ShouldDropKiss(p Packet) bool { return p.Stratum == 0 && p.ReferenceID == 0 }
