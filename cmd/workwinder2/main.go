package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/slynickel/workwinder2/timer"
)

func main() {
	basic2()
}

const refreshInterval = 500 * time.Millisecond

var (
	table         *tview.Table
	app           *tview.Application
	activeRow     *int
	overallTimers []*timer.Timer
)

func basic2() {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	zzz := 0
	activeRow = &zzz

	fakenames := []string{"INTERNAL: STOP TIMER", "Internal (4)", "Mgmt (5)", "7091 Meetings"}
	for i, v := range fakenames {
		overallTimers = append(overallTimers, timer.New(i, v))
	}

	headerRow := strings.Split("Stop/Start Time Label", " ")

	app = tview.NewApplication()
	table = tview.NewTable().SetBorders(false).SetSelectable(true, false).
		SetFixed(1, 1).SetSeparator(tview.Borders.Vertical)

	// set header
	for c := 0; c < len(headerRow); c++ {
		table.SetCell(0, c,
			tview.NewTableCell(headerRow[c]).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter))
	}

	for _, tmr := range overallTimers {
		tmr.InitVisuals(table)
	}

	table.Select(1, 0).SetFixed(1, 2).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	})

	table.SetSelectedFunc(func(newrow int, column int) {
		if *activeRow == newrow || newrow == 0 { // it shouldn't catch this case but in case it does
			return
		}
		// TODO handle stopped state and allow for stopping
		if *activeRow != 0 {
			overallTimers[*activeRow].Stop(table)
		}
		overallTimers[newrow].Start(table)
		table.SetSelectable(true, false)
		activeRow = &newrow
	})

	// Redraw always follows this
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		f.WriteString(fmt.Sprintf("key: %v\n", event.Key()))
		if event.Key() == tcell.KeyCtrlP {
			overallTimers = append(overallTimers, timer.New(len(overallTimers), "todo, set"))
			overallTimers[len(overallTimers)-1].InitVisuals(table)
		} else if event.Key() == tcell.KeyCtrlN {
			table.RemoveRow(overallTimers[len(overallTimers)-1].Index)
		}
		return event
	})

	go updateSelected()

	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func updateSelected() {
	for {
		time.Sleep(refreshInterval)
		if activeRow == nil {
			continue
		}
		if *activeRow == 0 { // stopped
			continue
		}
		app.QueueUpdateDraw(func() {
			c := table.GetCell(*activeRow, 1)
			c.SetText(overallTimers[*activeRow].CalculateVisibleDuration(time.Now()))
		})
	}
}

func formatDuration(input time.Duration) string {
	dur := input
	h := dur / time.Hour
	dur = dur - (h * time.Hour)
	m := dur / time.Minute
	dur = dur - (m * time.Minute)
	s := dur / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
