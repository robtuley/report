package report

import (
	"time"
)

func Global(payload Data) {
	addGlobalChannel <- payload
}

func init() {
	go func() {
		globals := Data{}

		for {
			select {
			case evt, more := <-rawEventChannel:
				if !more {
					drainChannel <- true
					return
				}
				evt["timestamp"] = time.Now().Unix()
				for k, v := range globals {
					evt[k] = v
				}
				withGlobalsEventChannel <- evt
			case g := <-addGlobalChannel:
				for k, v := range g {
					globals[k] = v
				}
			}
		}
	}()
}
