package main

import (
	"fmt"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	nstyle "github.com/aarzilli/nucular/style"
	log "github.com/sirupsen/logrus"
)

type settings struct {
	// Visualization options

	FlashOnStop bool
	Theme       iThemes

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

	textBoxes        []nucular.TextEditor
	textBoxMaxLength int
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

type iThemes int

const (
	DarkTheme iThemes = iota
	DefaultTheme
	RedTheme
	WhiteTheme
)

type nTheme struct {
	Name  string
	theme nstyle.Theme
}

var (
	zeroTime   time.Time
	themeNames = []string{"Dark", "Default", "Red", "White"}
	themesN    = []nstyle.Theme{nstyle.DarkTheme, nstyle.DefaultTheme, nstyle.RedTheme, nstyle.WhiteTheme}
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

// Timer is the backend structre of data for workwinder
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
	history []Event
	// isTrackingStoppedTime is true when this timer is tracking the amount of time no timer is active
	isTrackingStoppedTime bool
}

func newSettings() (s *settings) {
	s = &settings{}
	// TODO load settings file
	s.DefaultColumns = []string{"mgmt", "test1", "test2"}
	s.FlashOnStop = true
	s.Theme = DarkTheme
	s.textBoxMaxLength = 256

	// TODO load current active timer if relevant

	// Initializing structs
	s.stopTimer = NewTimer(st, true)
	s.textBoxes = make([]nucular.TextEditor, len(s.DefaultColumns))
	for i, colName := range s.DefaultColumns {
		s.timers = append(s.timers, NewTimer(colName, false))

		s.textBoxes[i].Buffer = []rune(s.timers[i].Name)
		s.textBoxes[i].Maxlen = s.textBoxMaxLength
		// TODO set s.textBoxes[*].EditFlags but unclear on which ones
	}
	// start with stopped timer running
	s.stateHandler(-1)
	return s
}

func (s *settings) run(w *nucular.Window) {
	now := time.Now()

	// Settings Drop Down
	if w.TreePush(nucular.TreeTab, "Settings", false) {
		w.Row(30).Dynamic(2)
		w.Label("Theme: ", "LC")
		oldTheme := s.Theme
		s.Theme = iThemes(w.ComboSimple(themeNames, int(s.Theme), 25))
		if oldTheme != s.Theme {
			w.Master().SetStyle(nstyle.FromTheme(themesN[s.Theme], w.Master().Style().Scaling))
		}
		w.TreePop()
	}

	// Header Rows
	w.Row(30).Static(38, 38, 280)
	if w.Button(label.ST(label.SymbolPlus, "", "RC"), false) {
		fmt.Print("up")
	}
	if w.Button(label.ST(label.SymbolMinus, "", "RC"), false) {
		fmt.Print("down")
	}
	if w.ButtonText("Report") {
		fmt.Print("report")
	}

	w.Row(30).Static(80, 280)
	if w.ButtonText(s.stopTimer.stopText()) {
		log.Debug("type: button, index: (stopped), line name: (stopped)")
		s.stateHandler(-1)
	}
	w.Label(s.CacluateCurrentOverallTotal(now), "CC")

	// Body Rows
	for i := range s.timers {
		// play/running, textbox, duration of timer
		w.Row(30).Static(80, 200, 80)
		if w.Button(s.timers[i].timerText(), false) {
			log.Debugf("type: button, index: %d, line name: %s", i, s.timers[i].Name)
			s.stateHandler(i)
		}
		s.textBoxes[i].Edit(w)
		w.Label(s.timers[i].GetCurrentTotal(now), "CC")
	}
}

func (s *settings) CacluateCurrentOverallTotal(now time.Time) string {
	prefix := "Total: "
	if s.activeIndex == -1 {
		return prefix + FormatDuration(s.overAllSavedTotal)
	}
	return prefix + FormatDuration(now.Sub(s.timers[s.activeIndex].StartTimestamp)+s.overAllSavedTotal)
}

// stateHandler should only be called on open and button presses
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
}

func (s *settings) addTimer() {

}

// NewTimer creates a new Timer instance with a state of TimerCreated
func NewTimer(name string, isTrackingStoppedTime bool) *Timer {
	now := time.Now()
	t := &Timer{
		Name:                  name,
		State:                 TimerCreate,
		isTrackingStoppedTime: isTrackingStoppedTime,
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
	t.history = append(t.history, Event{
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
	return label.ST(label.SymbolTriangleRight, "", "RC")
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
