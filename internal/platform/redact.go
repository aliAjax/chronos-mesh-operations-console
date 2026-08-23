package platform

import "strings"

func Redact(v string) string {
	if len(v) <= 4 {
		return "****"
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
func SafeAddress(addr string) string {
	if i := strings.LastIndex(addr, ":"); i > 0 {
		return addr[:i] + ":*"
	}
	return addr
}
func SecretFields(m map[string]string) map[string]string {
	o := map[string]string{}
	for k, v := range m {
		lk := strings.ToLower(k)
		if strings.Contains(lk, "key") || strings.Contains(lk, "secret") || strings.Contains(lk, "password") {
			o[k] = "[redacted]"
		} else {
			o[k] = v
		}
	}
	return o
}
