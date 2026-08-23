package metrics

import (
	"fmt"
	"sync/atomic"
)

type Registry struct {
	Packets  atomic.Uint64
	Requests atomic.Uint64
	Errors   atomic.Uint64
}

func (r *Registry) Text() string {
	return fmt.Sprintf("chronos_ntp_packets_total %d\nchronos_control_requests_total %d\nchronos_errors_total %d\n", r.Packets.Load(), r.Requests.Load(), r.Errors.Load())
}
