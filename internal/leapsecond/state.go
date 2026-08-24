package leapsecond

import (
	"sync"
	"time"
)

type State struct {
	mu      sync.RWMutex
	Name    string
	At      time.Time
	Message string
}

func (s *State) Set(name, msg string, at time.Time) {
	s.mu.Lock()
	s.Name = name
	s.Message = msg
	s.At = at
	s.mu.Unlock()
}
func (s *State) Get() (string, string, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Name, s.Message, s.At
}
