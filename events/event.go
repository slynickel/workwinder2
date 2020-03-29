package events

import "time"

// States
const (
	// Stop indicates the timer was stopped
	Stopped string = "Stopped"
	// Running indicates the timer is running
	Running string = "Running"
	// Create indicates the timer was created
	Created string = "Created"
	// Removed indicates the timer was removed
	Removed string = "Removed"
)

// Event is a change of state for a timer
type Event struct {
	// Timestamp when the event occured
	Timestamp time.Time
	// State the time moved to
	State string
	// TimerName when the event occurs
	TimerName string
	// Duration is the duration of the event, may be 0 if starting or removing
	Duration time.Duration
	// Total is the total duration at the time of the event
	Total time.Duration
}
