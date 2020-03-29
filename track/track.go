package track

import (
	"time"

	"github.com/slynickel/workwinder2/events"
)

type Work interface {
	Init()
	Start()
	Stop()
	Running() bool

	GetEvents() []events.Event

	Duration() time.Duration

	SetName(string)
	GetName() string
}
