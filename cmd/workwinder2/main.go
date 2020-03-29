package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	basic2()
}

const refreshInterval = 500 * time.Millisecond
const stopped = "▶"
const running = "Running"

var (
	table     *tview.Table
	app       *tview.Application
	activeRow *int

	overallTimers *[]Poc
)

type Poc struct {
	previous time.Time
	state    string
	name     string
	total    time.Duration
}

func basic2() {
	toss := 0
	activeRow = &toss
	overallTimers = &[]Poc{
		Poc{
			state:    running,
			name:     "INTERNAL: STOP TIMER",
			previous: time.Now(),
		},
		Poc{
			state:    stopped,
			name:     "Internal(4)",
			previous: time.Now(),
		},
		Poc{
			state:    stopped,
			name:     "Mgmt(5)",
			previous: time.Now(),
		},
		Poc{
			state:    stopped,
			name:     "7091 Meetings",
			previous: time.Now(),
		},
	}

	headerRow := strings.Split("Stop/Start Time Label", " ")

	app = tview.NewApplication()
	table = tview.NewTable().SetBorders(true)

	cols, rows := len(headerRow), len(*overallTimers) // the header row is the hidden stopped timer

	// set header
	for c := 0; c < cols; c++ {
		table.SetCell(0, c,
			tview.NewTableCell(headerRow[c]).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter))
	}
	// set body
	for r := 1; r < rows; r++ {
		for c := 0; c < cols; c++ {
			row := formatPOC((*overallTimers)[r])
			table.SetCell(r, c,
				tview.NewTableCell(row[c]).
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

	// Redraw always follows this
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlP { // PLUS
			table.InsertRow(4)
			cols := 3
			for c := 0; c < cols; c++ {
				row := strings.Split("example1 2 4", " ")
				table.SetCell(4, c,
					tview.NewTableCell(row[c]).
						SetTextColor(tcell.ColorWhite).
						SetAlign(tview.AlignCenter))
			}
		}
		if event.Key() == tcell.KeyCtrlM { // MINUS
			table.RemoveRow(4)
		}
		return event
	})

	table.SetSelectedFunc(func(newrow int, column int) {
		if *activeRow != 0 {
			// toggle the previous timer VISUAL state
			// toggleRowState(table, *activeRow)
			cell := table.GetCell(*activeRow, 0)
			if cell.Text == stopped {
				cell.SetText(running)
			} else {
				cell.SetText(stopped)
			}
			// toogle the previous timer STORAGE state, potentially update the name
			(*overallTimers)[*activeRow].togglePocState(table.GetCell(*activeRow, 2).Text)
		}

		// toogle the new timer VISUAL state
		//toggleRowState(table, newrow)
		// toogle the new timer STORAGE state, potentially udpate the name
		cell := table.GetCell(newrow, 0)
		if cell.Text == stopped {
			cell.SetText(running)
		} else {
			cell.SetText(stopped)
		}
		(*overallTimers)[newrow].togglePocState(table.GetCell(newrow, 2).Text)
		// table.SetSelectable(true, false)
		activeRow = &newrow
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
		if *activeRow == 0 {
			continue
		}
		app.QueueUpdateDraw(func() {
			c := table.GetCell(*activeRow, 1)
			oldDuration := (*overallTimers)[*activeRow].total
			eventStart := (*overallTimers)[*activeRow].previous
			runningDuration := time.Now().Sub(eventStart)
			c.SetText(
				formatDuration((runningDuration + oldDuration).Round(time.Second)),
			)
		})
	}
}

func formatPOC(s Poc) []string {
	return []string{
		s.state,
		formatDuration(s.total),
		s.name,
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

func (p *Poc) togglePocState(cellName string) {
	if p.state == stopped {
		p.previous = time.Now()
		p.state = running
		return
	}

	if p.state == running {
		now := time.Now()
		dur := now.Sub(p.previous)
		p.total = p.total + dur
		p.state = stopped
		p.previous = now
		return
	}
}

func toggleRowState(t *tview.Table, rawRowIndex int) {
	cols := 3
	if rawRowIndex == 0 { // this would handle the stoped timer which needs to happen at somepoint
		return
	}

	for c := 0; c < cols; c++ {
		cell := table.GetCell(rawRowIndex, c)
		if cell.BackgroundColor == tcell.ColorDefault {
			cell.SetBackgroundColor(tcell.ColorWhite).SetTextColor(tcell.ColorBlack)
		} else {
			cell.SetBackgroundColor(tcell.ColorDefault).SetTextColor(tcell.ColorWhite)
		}
		if c == 0 {
			if cell.Text == stopped {
				cell.SetText(running)
			} else {
				cell.SetText(stopped)
			}
		}
	}
}
