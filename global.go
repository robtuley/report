package report

import (
	"time"
)

func Global(payload Data) {
	channel.AddGlobal <- payload
}

func init() {
	go func() {
		globals := Data{}

		for {
			select {
			case evt, more := <-channel.RawEvents:
				if !more {
					channel.Drain <- true
					return
				}
				evt["timestamp"] = time.Now().Unix()
				for k, v := range globals {
					evt[k] = v
				}
				channel.WithGlobals <- evt
			case g := <-channel.AddGlobal:
				for k, v := range g {
					globals[k] = v
				}
			}
		}
	}()
}
