package ntp

import (
	"context"
	"github.com/chronos-mesh/chronos-mesh/internal/platform"
	"net"
	"sync/atomic"
	"time"
)

type Model interface{ Snapshot() Snapshot }
type Snapshot struct {
	Now                       time.Time
	Stratum                   uint8
	Reference                 time.Time
	Synchronized              bool
	RootDelay, RootDispersion time.Duration
}
type Server struct {
	addr    string
	model   Model
	log     Logger
	conn    *net.UDPConn
	packets atomic.Uint64
}
type Logger interface {
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

func NewServer(addr string, m Model, l Logger) *Server { return &Server{addr: addr, model: m, log: l} }
func (s *Server) Start(ctx context.Context) error {
	a, e := net.ResolveUDPAddr("udp", s.addr)
	if e != nil {
		return e
	}
	s.conn, e = net.ListenUDP("udp", a)
	if e != nil {
		return e
	}
	go func() { <-ctx.Done(); s.conn.Close() }()
	buf := make([]byte, 2048)
	for {
		n, peer, e := s.conn.ReadFromUDP(buf)
		if e != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return e
			}
		}
		req, e := Decode(buf[:n])
		if e != nil || req.Mode != 3 {
			continue
		}
		snap := s.model.Snapshot()
		now := snap.Now
		resp := RequestResponse(req, now, snap.Stratum, snap.Reference)
		if !snap.Synchronized {
			resp.LI = 3
		}
		resp.RootDelay = uint32(EncodeDuration(snap.RootDelay))
		resp.RootDispersion = uint32(EncodeDuration(snap.RootDispersion))
		out, _ := Encode(resp)
		_, _ = s.conn.WriteToUDP(out, peer)
		s.packets.Add(1)
	}
}
func (s *Server) Packets() uint64 { return s.packets.Load() }
func (s *Server) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

var _ platform.Clock
