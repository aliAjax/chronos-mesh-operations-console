package source

import (
	"context"
	"github.com/chronos-mesh/chronos-mesh/internal/ntp"
	"net"
	"sync"
	"time"
)

type Sample struct {
	Name                 string
	Address              string
	Offset, Delay        time.Duration
	Reach                uint8
	Jitter, RootDistance time.Duration
	Last                 time.Time
	Accepted             bool
	Reason               string
	Stratum              uint8
}
type Provider interface {
	Poll(context.Context, string) (Sample, error)
}
type UDPProvider struct{ Timeout time.Duration }

func (p UDPProvider) Poll(ctx context.Context, addr string) (Sample, error) {
	s := Sample{Address: addr}
	c, e := net.DialTimeout("udp", addr, p.Timeout)
	if e != nil {
		return s, e
	}
	defer c.Close()
	now := time.Now().UTC()
	req := ntp.Packet{Version: 4, Mode: 3, Transmit: now}
	b, _ := ntp.Encode(req)
	if _, e = c.Write(b); e != nil {
		return s, e
	}
	_ = c.SetReadDeadline(time.Now().Add(p.Timeout))
	rb := make([]byte, 512)
	n, e := c.Read(rb)
	if e != nil {
		return s, e
	}
	resp, e := ntp.Decode(rb[:n])
	if e != nil {
		return s, e
	}
	recv := time.Now().UTC()
	s.Offset = ((resp.Receive.Sub(now)) + (resp.Transmit.Sub(recv))) / 2
	s.Delay = recv.Sub(now) - (resp.Transmit.Sub(resp.Receive))
	s.Stratum = resp.Stratum
	s.Last = recv
	s.Reach = 1
	return s, nil
}

type Manager struct {
	mu       sync.RWMutex
	samples  map[string]Sample
	provider Provider
	interval time.Duration
}

func NewManager(p Provider) *Manager {
	return &Manager{samples: map[string]Sample{}, provider: p, interval: time.Second}
}
func (m *Manager) Run(ctx context.Context, cfg []struct{ Name, Address string }) {
	t := time.NewTicker(m.interval)
	defer t.Stop()
	poll := func() {
		for _, c := range cfg {
			go func(c struct{ Name, Address string }) {
				s, e := m.provider.Poll(ctx, c.Address)
				s.Name = c.Name
				if e != nil {
					s.Reason = e.Error()
					s.Accepted = false
				} else {
					s.Accepted = true
					s.Reason = "candidate"
				}
				m.mu.Lock()
				m.samples[c.Name] = s
				m.mu.Unlock()
			}(c)
		}
	}
	poll()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			poll()
		}
	}
}
func (m *Manager) List() []Sample {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Sample, 0, len(m.samples))
	for _, s := range m.samples {
		out = append(out, s)
	}
	return out
}
