package metrics

import (
	"fmt"
	"strings"
)

func SanitizeLabel(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, "\\", "_")
	v = strings.ReplaceAll(v, "\"", "_")
	v = strings.ReplaceAll(v, "\n", "_")
	if len(v) > 64 {
		v = v[:64]
	}
	return v
}
func FormatGauge(name string, value float64) string { return name + " " + formatFloat(value) + "\n" }
func formatFloat(v float64) string {
	if v == 0 {
		return "0"
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", v), "0"), ".")
}
