package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var exampleToRun string

func init() {
	const (
		defaultGopher = "pocket"
		usage         = "the variety of gopher"
	)
	flag.StringVar(&exampleToRun, "option", defaultGopher, usage)
}

func main() {
	flag.Parse()
	switch exampleToRun {
	case "1":
		basic1()
	case "2":
		basic2()
	}
}

func basic1() {
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}

type WorkTimer interface {
	Start()
	Pause()
	IsActive()

	CurrentDuration() time.Duration

	SetName(string)
	GetName()
}

func (s *State) Start() {
	s.Active = true
}

type State struct {
	Duration  time.Duration
	Name      string
	Active    bool
	LastEvent time.Time
}

func basic2() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(true)
	lorem := strings.Split("Stop/Start Time Label s/t 00:00 example s/t 00:00 example2 s/t 00:00 example3", " ")
	cols, rows := 3, 4
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < 1 || r < 1 {
				color = tcell.ColorYellow
			}
			table.SetCell(r, c,
				tview.NewTableCell(lorem[word]).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			word = (word + 1) % len(lorem)
		}
	}
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		c := table.GetCell(row, column)
		c.SetTextColor(tcell.ColorRed)
		c.SetText(currentTimeString())
		table.SetSelectable(false, false)
		exampleWrite(row, column)
	})
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func exampleWrite(row int, column int) {
	content := []byte(fmt.Sprintf("time: %s, row: %d, column %d", time.Now(), row, column))
	err := ioutil.WriteFile("hello", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func currentTimeString() string {
	t := time.Now()
	return fmt.Sprintf(t.Format("Current time is 15:04:05"))
}
