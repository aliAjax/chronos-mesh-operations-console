package platform

import (
	"os"
	"strconv"
	"time"
)

func EnvString(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func EnvInt(k string, d int) int {
	v, _ := strconv.Atoi(os.Getenv(k))
	return v
}
func EnvDuration(k string, d time.Duration) time.Duration {
	v, e := time.ParseDuration(os.Getenv(k))
	if e != nil {
		return d
	}
	return v
}
