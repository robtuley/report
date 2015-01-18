package report

import (
	"encoding/json"
	"log"
)

func init() {
	go func() {
		for {
			data, more := <-withGlobalsEventChannel
			if !more {
				drainChannel <- true
				return
			}
			json, err := map2Json(data)
			if err != nil {
				log.Println("error:> json encoding ", data)
				continue
			}
			jsonEventChannel <- json
		}
	}()
}

func map2Json(d Data) (string, error) {
	jsonBytes, err := json.Marshal(d)
	return string(jsonBytes), err
}
