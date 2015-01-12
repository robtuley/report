package report

import (
	"encoding/json"
)

func init() {
	go serializeDataEvents(withGlobalsEventChannel, jsonEventChannel)
}

func serializeDataEvents(in chan Data, out chan string) {
	for {
		json, err := map2Json(<-in)
		if err == nil {
			out <- json
		}
	}
}

func map2Json(d Data) (string, error) {
	jsonBytes, err := json.Marshal(d)
	return string(jsonBytes), err
}
