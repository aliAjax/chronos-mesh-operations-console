package platform

import "fmt"

type ErrorKind string

const (
	ErrProtocol ErrorKind = "protocol"
	ErrConfig   ErrorKind = "config"
	ErrSource   ErrorKind = "source"
	ErrCrypto   ErrorKind = "crypto"
)

type Error struct {
	Kind ErrorKind
	Op   string
	Err  error
}

func (e *Error) Error() string { return fmt.Sprintf("%s %s: %v", e.Kind, e.Op, e.Err) }
func (e *Error) Unwrap() error { return e.Err }
func Wrap(k ErrorKind, op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s %s: %v", k, op, err)
}
