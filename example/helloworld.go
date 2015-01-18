package main

import (
	"github.com/robtuley/report"
	"io"
	"net/http"
)

func main() {
	defer report.Drain()
	report.StdOut()
	// OR more likely:
	// report.SplunkStorm("yourUrl", "yourProjectId", "yourAccessKey")

	// add data for all log events if mixed aggregation
	report.Global(report.Data{"application": "myAppName"})

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		// timer to record response time and details
		defer report.Tock(report.Tick(), "http.response", report.Data{
			"host": req.URL.Host,
			"path": req.URL.Path,
			"ua":   req.UserAgent(),
		})
		io.WriteString(res, "Hello World")
	})

	// an *info* event to provide context
	report.Info("http.listening", report.Data{"port": 8080})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		// an *actionable* event that needs resolution
		report.Action("http.listening.fail", report.Data{"error": err.Error()})
	}
}
