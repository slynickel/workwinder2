package timers

import (
	"sync"
	"time"
)

// TimerState represents states of a Chronometer
type TimerState int

// Possible states of a Chronometer
// If this changes the schema version should be bumped
const (
	// Stopped can represent a Chronometer being stopped and should also represent a Chronometer being created >> stopped
	Stopped TimerState = iota
	// Started is called when a StopWatch is started
	Started
	// Created should only be called when a StopWatch is originally created (not on loads)
	Created
	// Loaded is when a StopWatch is loaded
	Loaded
)

// Event is an event that can occur to a Chronometer
type Event struct {
	// Timestamp when the event occurred
	Timestamp time.Time
	// RunDuration is the length of time if the state was Started to Stopped. Zero otherwise
	RunDuration time.Duration

	// Name when the event occurs
	Name string
	// Total is the total duration at the Stopwatch at the time of the event of the event
	Total time.Duration
	// Status the time moved into
	State TimerState
}

// StopWatch is the backend structure of data for workwinder
type StopWatch struct {
	Name string
	// Total duration of events that have occurred on the Chronometer
	Total time.Duration
	// Status is the current state of the timer
	State TimerState
	// StartTimestamp is the time stamp set ONLY with the TimerStart event and cleared (set to 0) when stopped
	StartTimestamp time.Time
	// History of events
	History []Event
}

// Chronometer satisfies Timers interface and provides a way to interact with a number of
// different, effectively, stop watches. All stop watches within Chronometers can be stopped.
// The locking is unnecessary as channels aren't used but it's good practice.
type Chronometer struct {
	mux           sync.Mutex
	SchemaVersion string

	StopTimer *StopWatch
	TimerList []*StopWatch
}

// NewChronometer returns a Chronometer struct with a StopWatch tracking stopped
// time in a running state
func NewChronometer() Timers {
	c := &Chronometer{
		StopTimer:     newStopWatch("tracking_stopped"),
		SchemaVersion: "1",
	}
	c.StopTimer.stateChange(Started, time.Now())
	return c
}

// Timers is a list of operations to modify timers
type Timers interface {
	Add(name string) int
	Remove(force bool) error
	Start(index int) error
	Stop()

	UpdateName(index int, name string) error

	Names() []string
	Totals(now time.Time) []string
	States() []TimerState
	Total(now time.Time) string
	StopTotal(now time.Time) string
}
