package report

import (
	"time"
)

// Log informational event that will provide context to any events requiring action.
//
//     report.Info("http.request", report.Data{"path":req.URL.Path, "ua":req.UserAgent()})
//
func Info(event string, payload Data) {
	payload["timestamp"] = time.Now().Unix()
	payload["type"] = "info"
	payload["event"] = event

	rawEventChannel <- payload
}
