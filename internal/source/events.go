package source

import (
	"sync"
	"time"
)

type Event struct {
	ID     uint64
	Name   string
	State  string
	At     time.Time
	Reason string
}
type EventLog struct {
	mu    sync.Mutex
	next  uint64
	items []Event
	limit int
}

func NewEventLog(limit int) *EventLog { return &EventLog{limit: limit} }
func (e *EventLog) Append(name, state, reason string) Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.next++
	x := Event{ID: e.next, Name: name, State: state, Reason: reason, At: time.Now().UTC()}
	e.items = append(e.items, x)
	if len(e.items) > e.limit {
		e.items = e.items[len(e.items)-e.limit:]
	}
	return x
}
func (e *EventLog) After(id uint64) []Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	o := []Event{}
	for _, x := range e.items {
		if x.ID > id {
			o = append(o, x)
		}
	}
	return o
}
