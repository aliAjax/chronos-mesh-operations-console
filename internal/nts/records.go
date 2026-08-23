package nts

import (
	"encoding/binary"
	"fmt"
)

const (
	RecordNextProtocol uint16 = 0x0104
	RecordAEAD         uint16 = 0x0105
	RecordServer       uint16 = 0x0106
	RecordPort         uint16 = 0x0107
	RecordEnd          uint16 = 0x0100
)

type Record struct {
	Type     uint16
	Critical bool
	Body     []byte
}

func EncodeRecords(rs []Record) []byte {
	b := []byte{}
	for _, r := range rs {
		n := len(r.Body)
		typ := r.Type
		if r.Critical {
			typ |= 0x8000
		}
		x := make([]byte, 4+n)
		binary.BigEndian.PutUint16(x, typ)
		binary.BigEndian.PutUint16(x[2:], uint16(n))
		copy(x[4:], r.Body)
		b = append(b, x...)
	}
	return b
}
func DecodeRecords(b []byte) ([]Record, error) {
	o := []Record{}
	for off := 0; off < len(b); {
		if len(b)-off < 4 {
			return nil, fmt.Errorf("short NTS record")
		}
		typ := binary.BigEndian.Uint16(b[off:])
		n := int(binary.BigEndian.Uint16(b[off+2:]))
		if off+4+n > len(b) {
			return nil, fmt.Errorf("NTS record overflow")
		}
		o = append(o, Record{Type: typ & 0x7fff, Critical: typ&0x8000 != 0, Body: append([]byte(nil), b[off+4:off+4+n]...)})
		off += 4 + n
	}
	return o, nil
}
