package ntp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

var ErrMalformedExtension = errors.New("malformed NTP extension")

const PacketSize = 48

type Packet struct {
	LI                                   uint8
	Version                              uint8
	Mode                                 uint8
	Stratum                              uint8
	Poll                                 int8
	Precision                            int8
	RootDelay                            uint32
	RootDispersion                       uint32
	ReferenceID                          uint32
	Reference, Origin, Receive, Transmit time.Time
	Extensions                           []Extension
}
type Extension struct {
	Type     uint16
	Critical bool
	Data     []byte
}

func ntpTime(t time.Time) uint64 {
	if t.IsZero() {
		return 0
	}
	sec := uint64(t.Unix() + 2208988800)
	frac := uint64((uint64(t.Nanosecond()) << 32) / 1e9)
	return sec<<32 | frac
}
func fromNTP(v uint64) time.Time {
	if v == 0 {
		return time.Time{}
	}
	sec := int64(v>>32) - 2208988800
	ns := int64((v & 0xffffffff) * 1e9 >> 32)
	return time.Unix(sec, ns).UTC()
}
func Encode(p Packet) ([]byte, error) {
	b := make([]byte, PacketSize)
	b[0] = (p.LI << 6) | ((p.Version & 7) << 3) | (p.Mode & 7)
	b[1] = p.Stratum
	b[2] = byte(p.Poll)
	b[3] = byte(p.Precision)
	binary.BigEndian.PutUint32(b[4:], p.RootDelay)
	binary.BigEndian.PutUint32(b[8:], p.RootDispersion)
	binary.BigEndian.PutUint32(b[12:], p.ReferenceID)
	binary.BigEndian.PutUint64(b[16:], ntpTime(p.Reference))
	binary.BigEndian.PutUint64(b[24:], ntpTime(p.Origin))
	binary.BigEndian.PutUint64(b[32:], ntpTime(p.Receive))
	binary.BigEndian.PutUint64(b[40:], ntpTime(p.Transmit))
	for _, x := range p.Extensions {
		n := len(x.Data) + 4
		pad := (n + 3) &^ 3
		old := len(b)
		b = append(b, make([]byte, pad)...)
		binary.BigEndian.PutUint16(b[old:], x.Type)
		binary.BigEndian.PutUint16(b[old+2:], uint16(n))
		copy(b[old+4:], x.Data)
	}
	return b, nil
}
func Decode(b []byte) (Packet, error) {
	if len(b) < PacketSize {
		return Packet{}, fmt.Errorf("short NTP packet: %d", len(b))
	}
	p := Packet{LI: b[0] >> 6, Version: (b[0] >> 3) & 7, Mode: b[0] & 7, Stratum: b[1], Poll: int8(b[2]), Precision: int8(b[3]), RootDelay: binary.BigEndian.Uint32(b[4:]), RootDispersion: binary.BigEndian.Uint32(b[8:]), ReferenceID: binary.BigEndian.Uint32(b[12:]), Reference: fromNTP(binary.BigEndian.Uint64(b[16:])), Origin: fromNTP(binary.BigEndian.Uint64(b[24:])), Receive: fromNTP(binary.BigEndian.Uint64(b[32:])), Transmit: fromNTP(binary.BigEndian.Uint64(b[40:]))}
	for off := 48; off < len(b); {
		if len(b)-off < 4 {
			return Packet{}, fmt.Errorf("truncated extension")
		}
		typ := binary.BigEndian.Uint16(b[off:])
		n := int(binary.BigEndian.Uint16(b[off+2:]))
		if n < 4 || off+n > len(b) {
			return Packet{}, fmt.Errorf("invalid extension length")
		}
		p.Extensions = append(p.Extensions, Extension{Type: typ, Data: append([]byte(nil), b[off+4:off+n]...)})
		off += (n + 3) &^ 3
	}
	return p, nil
}
func RequestResponse(req Packet, now time.Time, stratum uint8, ref time.Time) Packet {
	return Packet{LI: 0, Version: req.Version, Mode: 4, Stratum: stratum, Poll: req.Poll, Precision: -20, Reference: ref, Origin: req.Transmit, Receive: now, Transmit: now}
}
