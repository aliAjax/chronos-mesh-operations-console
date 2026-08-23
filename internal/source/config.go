package source

import "time"

type Config struct {
	Name, Address string
	Poll          time.Duration
	MaxOffset     time.Duration
	Enabled       bool
}

func (c Config) Normalize() Config {
	if c.Poll <= 0 {
		c.Poll = time.Second
	}
	if c.MaxOffset <= 0 {
		c.MaxOffset = 2 * time.Second
	}
	return c
}
