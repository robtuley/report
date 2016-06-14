// Example daemon demonstrating testable via report package logging
package main

import (
	"time"

	"github.com/robtuley/report"
)

func main() {
	defer report.Drain()
	report.StdOut()
	report.Global(report.Data{"service": "ticker"})
	report.RuntimeStatsEvery(time.Second * 10)

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
}
