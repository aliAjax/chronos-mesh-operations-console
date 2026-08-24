package listener

import (
	"context"
	"errors"
	"github.com/chronos-mesh/chronos-mesh/internal/ratelimit"
	"net"
	"sync"
	"time"
)

type UDP struct {
	Conn    *net.UDPConn
	Workers int
	Limit   *ratelimit.Limiter
	Handler func([]byte, *net.UDPAddr) []byte
	wg      sync.WaitGroup
}

func NewUDP(addr string, workers int, h func([]byte, *net.UDPAddr) []byte) (*UDP, error) {
	a, e := net.ResolveUDPAddr("udp", addr)
	if e != nil {
		return nil, e
	}
	c, e := net.ListenUDP("udp", a)
	if e != nil {
		return nil, e
	}
	u := &UDP{Conn: c, Workers: workers, Handler: h}
	return u, nil
}
func (u *UDP) Run(ctx context.Context) {
	u.wg.Add(1)
	defer u.wg.Done()
	buf := make([]byte, 2048)
	for {
		u.Conn.SetReadDeadline(timeNow())
		n, p, e := u.Conn.ReadFromUDP(buf)
		if e != nil {
			if errors.Is(e, net.ErrClosed) {
				return
			}
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		if u.Handler != nil {
			out := u.Handler(append([]byte(nil), buf[:n]...), p)
			if len(out) > 0 {
				u.Conn.WriteToUDP(out, p)
			}
		}
	}
}
func (u *UDP) Close()    { _ = u.Conn.Close(); u.wg.Wait() }
func timeNow() time.Time { return time.Now().Add(time.Second) }
