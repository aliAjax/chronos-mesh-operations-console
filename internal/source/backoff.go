package source

import "time"

type Backoff struct {
	Base, Max time.Duration
	attempt   int
}

func (b *Backoff) Next() time.Duration {
	d := b.Base
	for i := 0; i < b.attempt; i++ {
		d *= 2
		if d >= b.Max {
			d = b.Max
			break
		}
	}
	b.attempt++
	if d > b.Max {
		d = b.Max
	}
	return d
}
func (b *Backoff) Reset()      { b.attempt = 0 }
func (b Backoff) Attempt() int { return b.attempt }
