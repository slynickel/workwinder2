package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
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
	Wnd := nucular.NewMasterWindowSize(
		nucular.WindowTitle, // Most of the windows flags don't appear to do anything
		//nucular.WindowClosable|nucular.WindowMovable| << one these settings causes a crash?!
		"WorkWinder2",
		image.Point{X: 380, Y: 600},
		s.run,
	)
	Wnd.SetStyle(nstyle.FromTheme(s.Theme, *s.WindowScaling))

	go func() {
		for {
			time.Sleep(refreshInterval)
			if Wnd.Closed() {
				fmt.Printf("got here")
				break
			}
			Wnd.Changed()
		}
	}()
	Wnd.Main()
}

// TODO need to catch window closure and gracefully quit
