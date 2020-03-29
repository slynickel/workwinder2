package timer

import (
	"time"

	"github.com/slynickel/workwinder2/events"
)

// Timer represents the current and past state of an timer row
type Timer struct {
	// Current name of the time
	TimerName string
	// Total Duration of the timer
	Total time.Duration
	// Current state
	State string
	// LastEvent that occured to get into the current state
	LastEvent time.Time
	// History of events
	History []events.Event
}

func (t *Timer) Init() {
	t.State = events.Created
}

func (t *Timer) Start() {
	if t.State == events.Running {
		return
	}
	t.LastEvent = time.Now()
	t.State = events.Running
}

func (t *Timer) Stop() {
	if t.State == events.Stopped {
		return
	}
	now := time.Now()
	eventDuration := now.Sub(t.LastEvent)
	t.Total = t.Total + eventDuration
	stampEvent := events.Event{
		Timestamp: now,
		State:     events.Stopped,
		TimerName: t.TimerName,
		Duration:  eventDuration,
		Total:     t.Total,
	}
	t.History = append(t.History, stampEvent)
	t.State = events.Stopped
}

func (t *Timer) IsRunning() bool {
	if t.State != events.Running {
		return false
	}
	return true
}

func (t *Timer) Events() []events.Event {
	return t.History
}

func (t *Timer) Duration() time.Duration {
	return t.Total
}

func (t *Timer) SetName(name string) {
	t.TimerName = name
}

func (t *Timer) GetName() string {
	return t.TimerName
}
