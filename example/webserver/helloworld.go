// Example webserver demonstrating report package logging
package main

import (
	"github.com/rainchasers/report"
	"io"
	"net/http"
	"time"
)

func main() {
	// add data for all log events if mixed aggregation
	log := report.New(report.Data{"service": "myAppName"})

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		// timer to record response time and details
		defer log.Tock(log.Tick(), "http.response", report.Data{
			"host": req.URL.Host,
			"path": req.URL.Path,
			"ua":   req.UserAgent(),
		})
		io.WriteString(res, "Hello World")
	})

	// an *info* event to provide context
	log.Info("http.listening", report.Data{"port": 8080})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			// an *actionable* event that needs resolution, wait for log ack
			<-log.Action("http.listening.fail", report.Data{"error": err.Error()})
		}
	}()

	// to demo close exit after 60 seconds
	<-time.After(time.Second * 60)
}
