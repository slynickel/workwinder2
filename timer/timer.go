package timer

import (
	"fmt"
	"time"

	"github.com/slynickel/workwinder2/events"
)

// Timer represents the current and past state of an timer row
type Timer struct {
	// Current name of the time
	Name string
	// Total Duration of the timer
	Total time.Duration
	// Current state
	State string
	// LastEvent that occured to get into the current state
	LastEvent time.Time
	// History of events
	History []events.Event
}

func CreateTimer(name string) *Timer {
	now := time.Now()
	return &Timer{
		Name:      name,
		State:     events.Stopped,
		LastEvent: now,
		History: []events.Event{
			events.Event{
				Timestamp: now,
				State:     events.Created,
				TimerName: name,
			},
		},
	}
}

func (t *Timer) Start(name string) {
	t.Name = name
	if t.State == events.Running {
		return
	}
	t.LastEvent = time.Now()
	t.History = append(t.History, events.Event{
		Timestamp: t.LastEvent,
		State:     events.Running,
		TimerName: name,
		Total:     t.Total,
	})
	t.State = events.Running
}

func (t *Timer) Stop(name string) {
	t.Name = name
	if t.State == events.Stopped {
		return
	}
	now := time.Now()
	eventDuration := now.Sub(t.LastEvent)
	t.Total = t.Total + eventDuration
	t.History = append(t.History, events.Event{
		Timestamp: now,
		State:     events.Stopped,
		TimerName: t.Name,
		Duration:  eventDuration,
		Total:     t.Total,
	})
	t.State = events.Stopped
}

// FormatForCell returns a correctly formatted string array
//  it assumes total is the duration to display
func (t *Timer) FormatForCell() []string {
	return []string{t.State, FormatDuration(t.Total), t.Name}
}

// FormatDuration takes a standard duration and and formats it
//  into 00:00:00
func FormatDuration(d time.Duration) string {
	dur := d.Round(time.Second)
	h := dur / time.Hour
	dur = dur - (h * time.Hour)
	m := dur / time.Minute
	dur = dur - (m * time.Minute)
	s := dur / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (t *Timer) CalculateVisibleDuration(compare time.Time) string {
	runningDuration := compare.Sub(t.LastEvent)
	dur := t.Total + runningDuration
	return FormatDuration(dur)
}

//////////////////////
