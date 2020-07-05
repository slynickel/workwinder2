package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/slynickel/workwinder2/timers"

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

	c timers.Timers

	textBoxes []nucular.TextEditor
}

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

func New() (s *settings) {
	s = &settings{}
	// TODO load settings file
	s.DefaultColumns = []string{"mgmt", "test1", "test2"}
	s.FlashOnStop = true
	s.Theme = nstyle.DefaultTheme
	s.WindowScaling = &scalingDefault
	s.TextBoxMaxLength = 40
	s.c = timers.NewChronometer()

	// TODO load current active timer if relevant
	for _, n := range s.DefaultColumns {
		s.addLine(n)
	}
	return s
}

func (s *settings) addLine(name string) {
	s.c.Add(name)
	s.textBoxes = append(s.textBoxes, nucular.TextEditor{
		Buffer: []rune(name),
		Maxlen: s.TextBoxMaxLength,
	})
	// TODO set s.textBoxes[*].EditFlags but unclear on which ones
}

// Run is the control loop that runs on the clock cycle
func (s *settings) Run(w *nucular.Window) {
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
		w.Label("Version: 0.0.1", "LC") // TODO pull this from build flag

		if resetStyle {
			w.Master().SetStyle(nstyle.FromTheme(s.Theme, *s.WindowScaling))
		}

		w.TreePop()
	}

	// Header Rows
	// Buttons: Plus, Minus, Report
	w.Row(30).Ratio(0.1, 0.1, 0.8)
	if w.Button(label.ST(label.SymbolPlus, "", "RC"), false) {
		s.addLine("")
	}
	if w.Button(label.ST(label.SymbolMinus, "", "RC"), false) {
		fmt.Printf("TODO down")
	}
	if w.ButtonText("Report") {
		fmt.Print("TODO report")
	}

	// Button Stop/Stopped, Total Time
	w.Row(30).Ratio(0.2, 0.8)
	if w.ButtonText(s.stopText()) {
		log.Debug("type: button, index: (stopped), line name: (stopped)")
		s.c.Stop()
	}
	w.Label(s.c.Total(now), "CC")

	states := s.c.States()
	names := s.c.Names()
	totals := s.c.Totals(now)
	// Body Rows
	for i := range states {
		// Do the update name first to ensure the name update occurs before the state change
		if names[i] != string(s.textBoxes[i].Buffer) {
			s.c.UpdateName(i, string(s.textBoxes[i].Buffer))
		}

		// Button play/running, textbox, Label duration of timer
		w.Row(30).Ratio(0.2, 0.6, 0.2)
		if w.Button(timerText(states[i]), false) {
			log.Debugf("type: button, index: %d, line name: %s", i, names[i])
			s.c.Start(i)
		}
		s.textBoxes[i].Edit(w)
		w.Label(totals[i], "CC")
	}
}

func (s settings) writeSettings() {
	b, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Errorf("loadfile: %v", err)
	}
	err = ioutil.WriteFile("data/test-settings.json", b, 0644)
	if err != nil {
		log.Errorf("loadfile: %v", err)
	}
}

func (s *settings) stopText() string {
	// TimerStop really means that the stopped timer isn't runner
	// stopped means the stop timer is running (it's backwards)
	for _, s := range s.c.States() {
		if s == timers.Started {
			return "Stop"
		}
	}
	return "Stopped"
}

func timerText(state timers.TimerState) label.Label {
	if state == timers.Started {
		return label.T("Running")
	}
	return label.ST(label.SymbolTriangleRight, " ", "LC")
}
