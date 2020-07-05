package timers

import (
	"time"
)

// Names returns list of StopWatch names
func (c *Chronometer) Names() []string {
	names := []string{}
	for _, s := range c.TimerList {
		names = append(names, s.Name)
	}
	return names
}

// Totals returns list of StopWatch time totals including
// the time from now of the active timer (i.e. it will tick each time it's calculated)
func (c *Chronometer) Totals(now time.Time) []string {
	totals := []string{}
	for _, s := range c.TimerList {
		totals = append(totals, formatDuration(s.watchTotal(now)))
	}
	return totals
}

// States returns the TimerState of each StopWatch
func (c *Chronometer) States() []TimerState {
	ts := []TimerState{}
	for _, s := range c.TimerList {
		ts = append(ts, s.State)
	}
	return ts
}

// Total returns the overall total of StopWatches and includes the time from
// now of the actively running stop StopWatch
func (c *Chronometer) Total(now time.Time) string {
	var dur time.Duration
	for _, s := range c.TimerList {
		dur = dur + s.watchTotal(now)
	}
	return formatDuration(dur)
}

// StopTotal is the running total of the stopped StopWatch in Chronometer
// if actively running includes the time from now
func (c *Chronometer) StopTotal(now time.Time) string {
	return formatDuration(c.StopTimer.watchTotal(now))
}
