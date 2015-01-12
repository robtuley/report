package report

// Log informational event that will provide context to any events requiring action.
//
//     report.Info("http.request", report.Data{"path":req.URL.Path, "ua":req.UserAgent()})
//
func Action(event string, payload Data) {
	payload["type"] = "action"
	payload["event"] = event

	rawEventChannel <- payload
}
