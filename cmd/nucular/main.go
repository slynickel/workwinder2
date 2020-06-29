package main

import (
	"image"
	"os"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
	log "github.com/sirupsen/logrus"
)

// Setting much lower than 1000 results in updates stepping on each other leading to
// clock jitter between the over total and sub timers
const refreshInterval = 1000 * time.Millisecond

func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	s := newSettings()
	Wnd := nucular.NewMasterWindowSize(0, "WorkWinder2", image.Point{X: 360, Y: 600}, s.run)
	Wnd.SetStyle(style.FromTheme(themesN[s.Theme], 1.0))
	go func() {
		for {
			time.Sleep(refreshInterval)
			if Wnd.Closed() {
				break
			}
			Wnd.Changed()
		}
	}()
	Wnd.Main()
}
