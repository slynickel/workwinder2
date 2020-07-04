package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	nstyle "github.com/aarzilli/nucular/style"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mobile/event/mouse"
)

type settings struct {
	// Visualization options

	FlashOnStop      bool
	Theme            nstyle.Theme
	WindowScaling    *float64
	TextBoxMaxLength int

	// Timer Specific Options

	DefaultColumns []string

	// File path options

	SaveToFile   bool
	SaveFilePath string

	// Internal options

	overAllSavedTotal time.Duration
	activeIndex       int
	timers            []*Timer
	stopTimer         *Timer

	textBoxes []nucular.TextEditor
}

type TimerEvents int

const (
	// TimerStop indicates the timer was stopped
	TimerStop TimerEvents = iota
	// TimerStart indicates the timer was started
	TimerStart
	// TimerCreate indicates the timer was created
	TimerCreate
)

const (
	st = "trackingStoppedTime"
)

var (
	zeroTime time.Time
	// I think the labels on nstyle.Theme are silly. The Default should be
	// called Dark and vice versa
	themeNames = []string{"Dark", "White", "Red", "Default"}
	mapThemes  = map[int]nstyle.Theme{
		0: nstyle.DefaultTheme,
		1: nstyle.WhiteTheme,
		2: nstyle.RedTheme,
		3: nstyle.DarkTheme,
	}
	scalingDefault = 1.0
)

// Event is a change of state for a timer
type Event struct {
	// Timestamp when the event occured
	Timestamp time.Time
	// State the time moved into
	State TimerEvents
	// TimerName when the event occurs
	TimerName string
	// Total is the total duration at the time of the event
	Total time.Duration
}

// TODO rename Timer to Status or something like that so it isn't timer.Timer

// Timer is the backend structure of data for workwinder
type Timer struct {
	// Name of current time
	Name string
	// Total Duration of the timer
	Total time.Duration
	// State is the current state of the timer
	State TimerEvents
	// StartTimestamp is the time stamp set ONLY with the TimerStart event and cleared (set to 0) when stopped
	StartTimestamp time.Time
	// History of events
	History []Event
	// isTrackingStoppedTime is true when this timer is tracking the amount of time no timer is active
	IsTrackingStoppedTime bool
}

func newSettings() (s *settings) {
	s = &settings{}
	// TODO load settings file
	s.DefaultColumns = []string{"mgmt", "test1", "test2"}
	s.FlashOnStop = true
	s.Theme = nstyle.DefaultTheme
	s.WindowScaling = &scalingDefault
	s.TextBoxMaxLength = 40

	// TODO load current active timer if relevant

	// Initializing structs
	s.stopTimer = NewTimer(st, true)
	for _, colName := range s.DefaultColumns {
		s.addTimer(colName)
	}
	// start with stopped timer running
	s.stateHandler(-1)
	return s
}

func (s *settings) addTimer(name string) {
	index := len(s.timers)
	s.timers = append(s.timers, NewTimer(name, false))
	s.textBoxes = append(s.textBoxes, nucular.TextEditor{
		Buffer: []rune(s.timers[index].Name),
		Maxlen: s.TextBoxMaxLength,
	})
	// TODO set s.textBoxes[*].EditFlags but unclear on which ones
}

// run is the control loop that runs on the clock cycle
func (s *settings) run(w *nucular.Window) {
	now := time.Now()

	// Settings Drop Down
	if w.TreePush(nucular.TreeTab, "Settings", false) {
		resetStyle := false

		w.Row(30).Dynamic(2)
		w.Label("Theme: ", "LC")
		oldTheme := s.Theme
		s.Theme = mapThemes[w.ComboSimple(themeNames, int(s.Theme), 25)]
		if oldTheme != s.Theme {
			resetStyle = true
		}

		w.Row(30).Dynamic(3)
		w.Label("Window Scaling: ", "LC")
		w.Label(fmt.Sprintf("%.1f", *s.WindowScaling), "RC")

		w.SliderFloat(0.5, s.WindowScaling, 4, 0.1)
		// Only rescale if the mouse isn't pressed down and the value changed
		if w.Master().Style().Scaling != *s.WindowScaling && !w.Input().Mouse.Down(mouse.ButtonLeft) {
			resetStyle = true
		}

		w.Row(30).Dynamic(1)
		w.Label("Version: 0.0.1", "LC")

		if resetStyle {
			w.Master().SetStyle(nstyle.FromTheme(s.Theme, *s.WindowScaling))
		}

		w.TreePop()
	}

	// Header Rows
	// Buttons: Plus, Minus, Report
	w.Row(30).Ratio(0.1, 0.1, 0.8)
	if w.Button(label.ST(label.SymbolPlus, "", "RC"), false) {
		s.addTimer("")
	}
	if w.Button(label.ST(label.SymbolMinus, "", "RC"), false) {
		fmt.Printf("TODO down")
	}
	if w.ButtonText("Report") {
		fmt.Print("TODO report")
	}

	// Button Stop/Stopped, Total Time
	w.Row(30).Ratio(0.2, 0.8)
	if w.ButtonText(s.stopTimer.stopText()) {
		log.Debug("type: button, index: (stopped), line name: (stopped)")
		s.stateHandler(-1)
	}
	w.Label(s.CalculateCurrentOverallTotal(now), "CC")

	// Body Rows
	for i := range s.timers {
		// Button play/running, textbox, Label duration of timer
		w.Row(30).Ratio(0.2, 0.6, 0.2)
		if w.Button(s.timers[i].timerText(), false) {
			log.Debugf("type: button, index: %d, line name: %s", i, s.timers[i].Name)
			s.stateHandler(i)
		}
		s.textBoxes[i].Edit(w)
		w.Label(s.timers[i].GetCurrentTotal(now), "CC")
	}
}

func (s *settings) CalculateCurrentOverallTotal(now time.Time) string {
	prefix := "Total: "
	if s.activeIndex == -1 {
		return prefix + FormatDuration(s.overAllSavedTotal)
	}
	return prefix + FormatDuration(now.Sub(s.timers[s.activeIndex].StartTimestamp)+s.overAllSavedTotal)
}

// stateHandler should only be called on open and button presses otherwise you may end up writing to disk a lot
func (s *settings) stateHandler(indexOfActive int) {
	s.activeIndex = indexOfActive
	if indexOfActive == -1 {
		s.stopTimer.stateChange(TimerStart, st)
	} else {
		s.stopTimer.stateChange(TimerStop, st)
	}

	for i, t := range s.timers {
		if i == indexOfActive {
			t.stateChange(TimerStart, string(s.textBoxes[i].Buffer))
			continue
		}
		t.stateChange(TimerStop, string(s.textBoxes[i].Buffer))
	}
	s.overAllSavedTotal = 0
	for _, t := range s.timers {
		s.overAllSavedTotal = s.overAllSavedTotal + t.Total
	}
	s.writeData()
}

func (s settings) writeSettings() {
	b, err := json.Marshal(s)
	if err != nil {
		log.Errorf("loadfile: %v", err)
	}
	err = ioutil.WriteFile("data/test-settings.json", b, 0644)
	if err != nil {
		log.Errorf("loadfile: %v", err)
	}
}

func (s settings) writeData() {
	b, err := json.Marshal(s.timers)
	if err != nil {
		log.Errorf("loadfile: %v", err)
	}
	err = ioutil.WriteFile("data/test.json", b, 0644)
	if err != nil {
		log.Errorf("loadfile: %v", err)
	}
}

// NewTimer creates a new Timer instance with a state of TimerCreated
func NewTimer(name string, isTrackingStoppedTime bool) *Timer {
	now := time.Now()
	t := &Timer{
		Name:                  name,
		State:                 TimerCreate,
		IsTrackingStoppedTime: isTrackingStoppedTime,
	}
	t.fileEvent(now)
	return t
}

func (t *Timer) stateChange(event TimerEvents, name string) {
	// Could implement a ChangeName event but for now just update the name and don't
	// log an event
	t.Name = name
	if t.State == event {
		return
	}

	now := time.Now()

	switch event {
	case TimerCreate:
		panic("TimerCreate not implemented in stateChange")

	case TimerStart:
		t.StartTimestamp = now

	case TimerStop:
		if t.StartTimestamp.IsZero() {
			break // timer wasn't previously running
		}
		t.Total = t.Total + now.Sub(t.StartTimestamp)
		t.StartTimestamp = zeroTime
	}
	t.State = event
	t.fileEvent(now)
}

func (t *Timer) fileEvent(now time.Time) {
	t.History = append(t.History, Event{
		Timestamp: now,
		State:     t.State,
		TimerName: t.Name,
		Total:     t.Total,
	})
}

func (t *Timer) stopText() string {
	// TimerStop really means that the stopped timer isn't runner
	// stopped means the stop timer is running (it's backwards)
	if t.State == TimerStart {
		return "Stopped"
	}
	return "Stop"
}

func (t *Timer) timerText() label.Label {
	if t.State == TimerStart {
		return label.T("Running")
	}
	return label.ST(label.SymbolTriangleRight, " ", "LC")
	//return label.S(label.SymbolTriangleRight)
}

func FormatDuration(d time.Duration) string {
	dur := d.Round(time.Second)
	h := dur / time.Hour
	dur = dur - (h * time.Hour)
	m := dur / time.Minute
	dur = dur - (m * time.Minute)
	s := dur / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (t *Timer) GetCurrentTotal(now time.Time) string {
	if t.StartTimestamp.IsZero() {
		return FormatDuration(t.Total)
	}
	runningDuration := now.Sub(t.StartTimestamp)
	return FormatDuration(t.Total + runningDuration)
}
