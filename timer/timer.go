package timer

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/slynickel/workwinder2/events"
)

// Timer represents the current and past state of an timer row
type Timer struct {
	// Current name of the time
	Name string
	// Total Duration of the timer
	Total time.Duration
	// Current state
	State string
	// Index is the counter index
	Index int
	// LastEvent that occured to get into the current state
	LastEvent time.Time
	// History of events
	History []events.Event
}

// Bucket holds everything else
type Bucket struct {
	Tmrs  []*Timer
	Table *tview.Table
}

func InitBucket() *Bucket {
	return &Bucket{
		Tmrs: []*Timer{},
	}
}

func (b *Bucket) Add(name string) int {
	i := len(b.Tmrs)
	b.Tmrs = append(b.Tmrs, New(i, name))
	b.Tmrs[i].InitVisuals(b.Table)
	return i
}

// Start returns the previously active row after stopping it
// and starting the timer at the new index
func (b *Bucket) Start(newIndex int) int {
	a := b.ActiveRow()
	if a == newIndex {
		return a
	}
	if a != -1 {
		b.Tmrs[a].Stop(b.Table)
	}

	if newIndex == 0 {
		b.Stopped()
		return a
	}
	b.Tmrs[newIndex].Start(b.Table)
	return a
}

func (b *Bucket) Stopped() {
	b.Tmrs[0].Start(b.Table)
}

func (b *Bucket) ActiveRow() int {
	for i, t := range b.Tmrs {
		if t.State == events.Running {
			return i
		}
	}
	return -1 // shouldn't be reachable
}

func (b *Bucket) RefreshActive() {
	a := b.ActiveRow()
	b.Table.GetCell(a, 1).SetText(b.Tmrs[a].CalculateVisibleDuration(time.Now()))
}

func New(index int, name string) *Timer {
	now := time.Now()
	return &Timer{
		Name:      name,
		State:     events.Stopped,
		Index:     index,
		LastEvent: now,
		History: []events.Event{
			{
				Timestamp: now,
				State:     events.Created,
				TimerName: name,
				Index:     index,
			},
		},
	}
}

// For now these return the same thing, in the future
// the input should be an interface or something
// and you marshal it from a file
func Load(index int) *Timer {
	now := time.Now()
	return &Timer{
		Name:      "", // TODO LOAD
		State:     events.Stopped,
		Index:     index,
		LastEvent: now,
		History: []events.Event{
			{
				Timestamp: now,
				State:     events.Created,
				TimerName: "", // TODO LOAD,
				Index:     index,
			},
		},
	}
}

func (tmr *Timer) Delete(t *tview.Table) {
	// TODO
}

func (t *Timer) InitVisuals(table *tview.Table) {
	vals := t.FormatForCell()
	for col, text := range vals {
		newcell := tview.NewTableCell(text)
		t.TCellSetState(newcell)
		table.SetCell(t.Index, col, newcell)
	}
}

func (tmr *Timer) UpdateVisuals(t *tview.Table) {
	vals := tmr.FormatForCell()
	for col, text := range vals {
		cell := t.GetCell(tmr.Index, col)
		cell.SetText(text)
		tmr.TCellSetState(cell)
	}
}

func (tmr *Timer) TCellSetState(t *tview.TableCell) {
	t.SetAlign(tview.AlignCenter)
	switch tmr.State {
	case events.Running:
		t.SetTextColor(tcell.ColorBlack).
			SetBackgroundColor(tcell.ColorWhite)
	case events.Stopped:
		t.SetTextColor(tcell.ColorDefault).
			SetBackgroundColor(tcell.ColorDefault)
	}
}

func (tmr *Timer) Start(t *tview.Table) {
	tmr.Name = t.GetCell(tmr.Index, 2).Text
	if tmr.State == events.Running {
		return
	}
	tmr.LastEvent = time.Now()
	tmr.History = append(tmr.History, events.Event{
		Timestamp: tmr.LastEvent,
		State:     events.Running,
		TimerName: tmr.Name,
		Total:     tmr.Total,
	})
	tmr.State = events.Running
	tmr.UpdateVisuals(t)
}

func (tmr *Timer) Stop(t *tview.Table) {
	tmr.Name = t.GetCell(tmr.Index, 2).Text
	if tmr.State == events.Stopped {
		return
	}
	now := time.Now()
	eventDuration := now.Sub(tmr.LastEvent)
	tmr.Total = tmr.Total + eventDuration
	tmr.History = append(tmr.History, events.Event{
		Timestamp: now,
		State:     events.Stopped,
		TimerName: tmr.Name,
		Duration:  eventDuration,
		Total:     tmr.Total,
	})
	tmr.State = events.Stopped
	tmr.UpdateVisuals(t)
}

// FormatForCell returns a correctly formatted string array
//  it assumes total is the duration to display
func (t *Timer) FormatForCell() []string {
	return []string{t.State, FormatDuration(t.Total), t.Name}
}

// FormatDuration takes a standard duration and and formats it
//  into 00:00:00
func FormatDuration(d time.Duration) string {
	dur := d.Round(time.Second)
	h := dur / time.Hour
	dur = dur - (h * time.Hour)
	m := dur / time.Minute
	dur = dur - (m * time.Minute)
	s := dur / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (t *Timer) CalculateVisibleDuration(compare time.Time) string {
	runningDuration := compare.Sub(t.LastEvent)
	dur := t.Total + runningDuration
	return FormatDuration(dur)
}
