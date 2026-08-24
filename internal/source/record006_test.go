package source

import (
	"reflect"
	"testing"
	"time"
)

func TestMedianOffsetPreservesSampleOrderQ006(t *testing.T) {
	in := []Sample{{Name: "late", Offset: 2 * time.Second}, {Name: "early", Offset: 0}, {Name: "middle", Offset: time.Second}}
	want := append([]Sample(nil), in...)
	if got := MedianOffset(in); got != time.Second {
		t.Fatalf("median=%s", got)
	}
	if !reflect.DeepEqual(in, want) {
		t.Fatalf("median rewrote input: %#v", in)
	}
}
