package selection

import (
	"fmt"
	"github.com/chronos-mesh/chronos-mesh/internal/source"
	"strings"
)

func Report(d Decision) string {
	var b strings.Builder
	fmt.Fprintf(&b, "reason=%s trusted=%d offset=%s\n", d.Reason, d.Trusted, d.Offset)
	for _, s := range d.Sources {
		fmt.Fprintf(&b, "%s %s accepted=%t reason=%s\n", s.Name, s.Offset, s.Accepted, s.Reason)
	}
	return b.String()
}
func Accepted(d Decision) []source.Sample {
	o := []source.Sample{}
	for _, s := range d.Sources {
		if s.Accepted {
			o = append(o, s)
		}
	}
	return o
}
