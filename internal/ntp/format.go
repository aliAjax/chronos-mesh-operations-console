package ntp

import (
	"encoding/binary"
	"fmt"
	"time"
)

func WriteTimestamp(dst []byte, off int, t time.Time) error {
	if off < 0 || off+8 > len(dst) {
		return fmt.Errorf("timestamp offset out of bounds")
	}
	binary.BigEndian.PutUint64(dst[off:], ntpTime(t))
	return nil
}
func ReadTimestamp(src []byte, off int) (time.Time, error) {
	if off < 0 || off+8 > len(src) {
		return time.Time{}, fmt.Errorf("timestamp offset out of bounds")
	}
	return fromNTP(binary.BigEndian.Uint64(src[off:])), nil
}
func IsZeroTimestamp(t time.Time) bool                      { return t.IsZero() }
func RoundTrip(orig, recv, tx, dst time.Time) time.Duration { return dst.Sub(orig) - tx.Sub(recv) }
func Offset(orig, recv, tx, dst time.Time) time.Duration    { return (recv.Sub(orig) + tx.Sub(dst)) / 2 }
func ValidateExtensionSize(n int) error {
	if n < 4 {
		return fmt.Errorf("extension shorter than header")
	}
	if n%4 != 0 {
		return fmt.Errorf("extension not aligned")
	}
	if n > 1024 {
		return fmt.Errorf("extension too large")
	}
	return nil
}
