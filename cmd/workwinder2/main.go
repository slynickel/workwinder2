package main

import (
	"fmt"
	"log"
	"os"
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
	b *timer.Bucket
	//table         *tview.Table
	app *tview.Application
	// activeRow *int
	// overallTimers []*timer.Timer
)

func basic2() {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	//zzz := 0
	//activeRow = &zzz

	app = tview.NewApplication()
	b = timer.InitBucket()
	b.Table = tview.NewTable().SetBorders(false).SetSelectable(true, false).
		SetFixed(1, 1).SetSeparator(tview.Borders.Vertical)

	fakenames := []string{"Stopped", "Internal (4)", "Mgmt (5)", "7091 Meetings"}
	for _, v := range fakenames {
		b.Add(v)
		// overallTimers = append(overallTimers, timer.New(i, v))
	}

	// headerRow := strings.Split("Stop/Start Time Stopped", " ")

	// // set header
	// for c := 0; c < len(headerRow); c++ {
	// 	b.Table.SetCell(0, c,
	// 		tview.NewTableCell(headerRow[c]).
	// 			SetTextColor(tcell.ColorYellow).
	// 			SetAlign(tview.AlignCenter))
	// }

	// for _, tmr := range overallTimers {
	// 	tmr.InitVisuals(table)
	// }

	b.Table.Select(1, 0).SetFixed(1, 2).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			b.Table.SetSelectable(true, false)
		}
	})

	b.Table.SetSelectedFunc(func(newrow int, column int) {
		// if *activeRow == newrow || newrow == 0 { // it shouldn't catch this case but in case it does
		// 	return
		// }
		b.Start(newrow)
		// TODO handle stopped state and allow for stopping
		// if *activeRow != 0 {
		// 	overallTimers[*activeRow].Stop(table)
		// }
		// overallTimers[newrow].Start(table)
		b.Table.SetSelectable(true, false)
		// activeRow = &newrow
	})

	// Redraw always follows this
	b.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		f.WriteString(fmt.Sprintf("key: %v\n", event.Key()))
		if event.Key() == tcell.KeyCtrlP {
			b.Add("TODO")
		} else if event.Key() == tcell.KeyCtrlN {
			//b.Table.RemoveRow(overallTimers[len(overallTimers)-1].Index)
		}
		return event
	})

	b.Start(0)

	go updateSelected()

	if err := app.SetRoot(b.Table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func updateSelected() {
	for {
		time.Sleep(refreshInterval)
		//a := b.ActiveRow()
		// if a == 0 { // stopped
		// 	continue
		// }
		app.QueueUpdateDraw(func() {
			b.RefreshActive()
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
