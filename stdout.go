package report

import (
	"log"
)

func StdOut() {
	go stdoutWriter()
}

func stdoutWriter() {
	for {
		json, more := <-channel.JsonEncoded
		if !more {
			channel.Drain <- true
			return
		}
		log.Println("json:>", json)
	}
}
