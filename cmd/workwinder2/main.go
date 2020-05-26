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
	b   *timer.Bucket
	app *tview.Application
)

func basic2() {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	app = tview.NewApplication()
	b = timer.InitBucket()
	b.Table = tview.NewTable().SetBorders(false).SetSelectable(true, false).
		SetFixed(1, 1).SetSeparator(tview.Borders.Vertical)

	fakenames := []string{"Stopped", "Internal (4)", "Mgmt (5)", "7091 Meetings"}
	for _, v := range fakenames {
		b.Add(v)
	}

	b.Table.Select(1, 0).SetFixed(1, 2).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			b.Table.SetSelectable(true, false)
		}
	})

	b.Table.SetSelectedFunc(func(newrow int, column int) {
		b.Start(newrow)
		b.Table.SetSelectable(true, false)
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
