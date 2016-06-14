// Example daemon demonstrating testable via report package logging
package main

import (
	"os"
	"time"

	"github.com/robtuley/report"
)

func main() {
	// setup logging output
	report.StdOut()
	report.Global(report.Data{"service": "ticker"})
	report.RuntimeStatsEvery(time.Second * 10)

	// use observer to accumulate log messages
	logC := make(chan report.Data, 50)
	log2Channel := func(d report.Data) {
		time.Sleep(time.Second * 2)
		// ^ sleep to demo drain works correctly
		select {
		case logC <- d:
		default:
			// (a non-blocking publish)
		}
	}
	report.Observe(log2Channel)

	// ticker daemon execution
	ticker := time.NewTicker(time.Second)
	report.Info("timer.start", report.Data{})

	go func() {
		seq := 0
		for range ticker.C {
			report.Info("timer.tick", report.Data{"sequence": seq})
			seq = seq + 1
		}
	}()

	time.Sleep(time.Millisecond * 15500)
	ticker.Stop()
	report.Info("timer.stop", report.Data{})

	// validate execution from log interface
	report.Drain()
	close(logC)

	nTicks := 0
	for d := range logC {
		if d["event"] == "timer.tick" {
			nTicks = nTicks + 1
		}
	}
	if nTicks == 15 {
		os.Exit(0)
	} else {
		os.Stderr.Write([]byte("*** FAILED *** (ticks != 15)"))
		os.Exit(1)
	}
}
