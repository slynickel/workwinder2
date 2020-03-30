package main

import (
	"fmt"
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
	zzz := 0
	activeRow = &zzz

	fakenames := []string{"INTERNAL: STOP TIMER", "Internal (4)", "Mgmt (5)", "7091 Meetings"}
	for _, v := range fakenames {
		overallTimers = append(overallTimers, timer.CreateTimer(v))
	}

	headerRow := strings.Split("Stop/Start Time Label", " ")

	app = tview.NewApplication()
	table = tview.NewTable().SetBorders(true)

	rows := len(overallTimers) // the header row is the hidden stopped timer

	// set header
	for c := 0; c < len(headerRow); c++ {
		table.SetCell(0, c,
			tview.NewTableCell(headerRow[c]).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter))
	}
	// set body
	for r := 1; r < rows; r++ {
		body := overallTimers[r].FormatForCell()
		for c := 0; c < len(body); c++ {
			table.SetCell(r, c,
				tview.NewTableCell(body[c]).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignCenter))
		}
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
		// also this is really gross
		if *activeRow != 0 {
			overallTimers[*activeRow].Stop(table.GetCell(*activeRow, 2).Text)
			table.GetCell(*activeRow, 0).SetText(overallTimers[*activeRow].State)
		}
		overallTimers[newrow].Start(table.GetCell(newrow, 2).Text)
		table.GetCell(newrow, 0).SetText(overallTimers[newrow].State)

		table.SetSelectable(true, false)
		activeRow = &newrow
	})

	// Redraw always follows this
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyCtrlP:
			rowToInsert := table.GetRowCount()
			table.InsertRow(rowToInsert)

			overallTimers = append(overallTimers, timer.CreateTimer("newline"))
			body := (overallTimers)[rowToInsert].FormatForCell()

			for c := 0; c < len(body); c++ {
				table.SetCell(rowToInsert, c,
					tview.NewTableCell(body[c]).
						SetTextColor(tcell.ColorWhite).
						SetAlign(tview.AlignCenter))
			}
			return nil
			// cant for the life of me get this to work
			// case tcell.KeyCtrlM:
			// 	table.RemoveRow(table.GetRowCount())
			// 	return nil
		}

		return event
	})

	go updateSelected()
	if err := app.SetRoot(table, true).Run(); err != nil {
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
