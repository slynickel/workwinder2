package timers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

// IndexOutOfRangeError satisfies the Error interface
type IndexOutOfRangeError struct {
	Index  int
	Length int
}

// TimerIsNonZero is an error that satisfies the Error interface and is used to flag
// the removal a StopWatch with non zero time
type TimerIsNonZero struct {
	Name  string
	Total string
	State string
}

func (e *IndexOutOfRangeError) Error() string {
	return fmt.Sprintf("Index: %d is out of range, length is %d", e.Index, e.Length)
}

func (e *TimerIsNonZero) Error() string {
	return fmt.Sprintf("Timer %s has time logged to it. It has a total of %s and is in state %s. If removed the will be removed from the total.",
		e.Name, e.Total, e.State)
}

func newStopWatch(name string) *StopWatch {
	now := time.Now()
	s := &StopWatch{ // Initial state is stopped but state change will be logged first
		Name: name,
	}
	s.stateChange(Created, now)
	s.stateChange(Stopped, now)
	return s
}

func (s *StopWatch) stateChange(state TimerState, now time.Time) {
	if s.State == state {
		return
	}
	var dur time.Duration

	switch state {
	case Created:
		// nothing needs to be adjusted
	case Started:
		s.StartTimestamp = now
	case Stopped:
		if s.StartTimestamp.IsZero() {
			break // timer wasn't previously running (created)
		}
		dur = now.Sub(s.StartTimestamp)
		s.Total = s.Total + dur
		s.StartTimestamp = time.Time{} // set to zero
	}

	s.State = state
	s.History = append(s.History, Event{
		Timestamp:   now,
		RunDuration: dur,
		Name:        s.Name,
		Total:       s.Total,
		State:       s.State,
	})
}

func (c *Chronometer) checkIndex(index int) error {
	if len(c.TimerList) < index || index < 0 {
		return &IndexOutOfRangeError{
			Index:  index,
			Length: len(c.TimerList),
		}
	}
	return nil
}

func (s *StopWatch) watchTotal(now time.Time) time.Duration {
	if s.StartTimestamp.IsZero() {
		return s.Total
	}
	return s.Total + now.Sub(s.StartTimestamp)
}

func formatDuration(d time.Duration) string {
	dur := d.Round(time.Second)
	h := dur / time.Hour
	dur = dur - (h * time.Hour)
	m := dur / time.Minute
	dur = dur - (m * time.Minute)
	s := dur / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (c *Chronometer) writeData() {
	b, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		log.Errorf("cannot marshal %v", err)
		return
	}
	err = ioutil.WriteFile("data/test.json", b, 0644)
	if err != nil {
		log.Errorf("on writefile: %v", err)
	}
}
