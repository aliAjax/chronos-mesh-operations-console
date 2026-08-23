package leapsecond

import (
	"fmt"
	"sort"
	"time"
)

func ValidateEntries(es []Entry) error {
	cp := append([]Entry(nil), es...)
	sort.Slice(cp, func(i, j int) bool { return cp[i].At.Before(cp[j].At) })
	for i, e := range cp {
		if e.Delta != 1 && e.Delta != -1 {
			return fmt.Errorf("entry %d invalid delta", i)
		}
		if i > 0 && e.At.Equal(cp[i-1].At) {
			return fmt.Errorf("duplicate leap instant")
		}
	}
	return nil
}
func Next(now time.Time, es []Entry) (Entry, bool) {
	for _, e := range es {
		if e.At.After(now) {
			return e, true
		}
	}
	return Entry{}, false
}
func Effective(now time.Time, es []Entry) int8 {
	var d int8
	for _, e := range es {
		if !e.At.After(now) {
			d += e.Delta
		}
	}
	return d
}
type Phase string
const ( PhaseEmpty Phase = "empty"; PhaseAnnounced Phase = "announced"; PhaseEffective Phase = "effective"; PhaseExpired Phase = "expired" )
func ResolvePhase(current Phase, now, expires time.Time, entries []Entry) Phase {
	if now.Before(expires) { if len(entries) > 0 { return PhaseAnnounced }; return PhaseEmpty }
	if len(entries) > 0 { return PhaseEffective }
	return PhaseExpired
}
