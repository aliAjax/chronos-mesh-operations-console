package configuration

import (
	"fmt"
	"strings"
)

func Validate(c Config) error {
	if c.Policies == nil {
		c.Policies["default"] = true
	}
	if strings.TrimSpace(c.NTPAddr) == "" || strings.TrimSpace(c.HTTPAddr) == "" {
		return fmt.Errorf("listen addresses required")
	}
	if c.RatePerSecond < 1 {
		return fmt.Errorf("rate_per_second must be positive")
	}
	if len(c.CookieKey) < 16 {
		return fmt.Errorf("cookie_key too short")
	}
	return nil
}
