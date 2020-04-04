package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slynickel/workwinder2/timer"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const refreshInterval = 500 * time.Millisecond

var (
	table     *tview.Table
	app       *tview.Application
	activeRow *int
	rows      []*TimerRow
)

func main() {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	app = tview.NewApplication()
	table = tview.NewTable().SetBorders(false).SetSelectable(false, false).
		SetFixed(1, 1).SetSeparator(tview.Borders.Vertical).
		SetSelectedStyle(tcell.ColorBlack, tcell.ColorBlue, 0)

	fakenames := []string{"INTERNAL: STOP TIMER", "Internal (4)", "Mgmt (5)", "7091 Meetings"}
	zero := 0
	for i, v := range fakenames {
		rows = append(rows, NewRow(i, v))
		rows[i].InitVisuals(table)
		if i == 0 {
			rows[i].T.Start(v)
			activeRow = &zero
			rows[i].UpdateVisuals(table)
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
		if *activeRow == newrow { // it shouldn't catch this case but in case it does
			return
		}

		newText := table.GetCell(*activeRow, 2).Text
		rows[*activeRow].T.Stop(newText)
		rows[*activeRow].UpdateVisuals(table)

		newText = table.GetCell(newrow, 2).Text
		rows[newrow].T.Start(newText)
		rows[newrow].UpdateVisuals(table)

		table.SetSelectable(true, false)
		activeRow = &newrow
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		overallCount := table.GetRowCount()

		f.WriteString(fmt.Sprintf("key: %v\n", event.Key()))

		if event.Key() == tcell.KeyCtrlP {
			rows = append(rows, NewRow(overallCount, "new row"))
			rows[overallCount].InitVisuals(table)
		} else if event.Key() == tcell.KeyCtrlN {
			// todo do a copy with one less row and remove the row
			// 			rows = append(rows, NewRow(i, "new row"))
			table.RemoveRow(overallCount - 1)
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
		app.QueueUpdateDraw(func() {
			c := table.GetCell(*activeRow, 1)
			c.SetText(rows[*activeRow].T.CalculateVisibleDuration(time.Now()))
		})
	}
}

type TimerRow struct {
	T     *timer.Timer
	index int
}
