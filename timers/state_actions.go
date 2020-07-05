package timers

import (
	"time"
)

// Add StopWatch to Chronometer and return the index (index 0) of the added StopWatch
func (c *Chronometer) Add(name string) int {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.TimerList = append(c.TimerList, newStopWatch(name))
	c.writeData()
	return len(c.TimerList) - 1
}

// Remove the last StopWatch from the Chronometer
// if the StopWatch to be removed has a non-zero total and force is not true
// the timer will not be removed and Error
func (c *Chronometer) Remove(force bool) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	panic("Chronometer.Remove is not implemented")
	return nil
}

func (c *Chronometer) Start(index int) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	now := time.Now()

	if err := c.checkIndex(index); err != nil {
		return err
	}

	if c.TimerList[index].State == Started { // A request was made to start an already running timer
		return nil
	}

	c.StopTimer.stateChange(Stopped, now)
	for i, s := range c.TimerList {
		if i == index {
			s.stateChange(Started, now)
			continue
		}
		s.stateChange(Stopped, now)
	}
	c.writeData()
	return nil
}

func (c *Chronometer) Stop() {
	c.mux.Lock()
	defer c.mux.Unlock()
	now := time.Now()

	if c.StopTimer.State == Started {
		return
	}

	c.StopTimer.stateChange(Started, now)
	for _, s := range c.TimerList {
		s.stateChange(Stopped, now)
	}
	c.writeData()
}

func (c *Chronometer) UpdateName(index int, name string) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if err := c.checkIndex(index); err != nil {
		return err
	}

	// TODO log this as an event?
	c.TimerList[index].Name = name
	return nil
}
