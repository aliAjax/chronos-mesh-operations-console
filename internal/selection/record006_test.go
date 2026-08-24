package selection

import (
	"github.com/chronos-mesh/chronos-mesh/internal/source"
	"reflect"
	"testing"
	"time"
)

func TestSelectPreservesCallerSamplesS006(t *testing.T) {
	in := []source.Sample{{Name: "late", Offset: time.Second, Delay: time.Millisecond, Reach: 1, Stratum: 2}, {Name: "early", Offset: 0, Delay: time.Millisecond, Reach: 1, Stratum: 2}}
	want := append([]source.Sample(nil), in...)
	_ = Select(in)
	if !reflect.DeepEqual(in, want) {
		t.Fatalf("caller samples rewritten: %#v", in)
	}
}

func TestWeightsRemainStableAcrossCallsS006(t *testing.T) {
	first := Weights([]source.Sample{{Name: "alpha", Delay: time.Millisecond}})
	_ = Weights([]source.Sample{{Name: "beta", Delay: time.Second}})
	if first[0].Name != "alpha" {
		t.Fatalf("weight snapshot reused: %#v", first)
	}
}

func TestAcceptedDoesNotRewriteDecisionS006(t *testing.T) {
	d := Decision{Sources: []source.Sample{{Name: "rejected", Accepted: false}, {Name: "accepted", Accepted: true}}}
	want := append([]source.Sample(nil), d.Sources...)
	got := Accepted(d)
	if len(got) != 1 || got[0].Name != "accepted" || !reflect.DeepEqual(d.Sources, want) {
		t.Fatalf("decision rewritten: %#v %#v", got, d.Sources)
	}
}
