// Example daemon demonstrating testable via report package logging
package main

import (
	"github.com/rainchasers/report"
	"os"
	"time"
)

func main() {
	// setup logging output
	log := report.New(report.Data{"service": "ticker"})

	// ticker daemon execution
	ticker := time.NewTicker(time.Second)
	log.Info("timer.start", report.Data{})

	go func() {
		seq := 0
		for range ticker.C {
			<-log.Info("timer.tick", report.Data{"sequence": seq})
			seq = seq + 1
		}
	}()

	time.Sleep(time.Millisecond * 15500)
	ticker.Stop()
	log.Info("timer.stop", report.Data{})

	// validate
	if log.Count("timer.tick") != 15 {
		os.Stderr.Write([]byte("*** FAILED *** (ticks != 15)"))
		os.Exit(1)
	}
}
