package report

import (
	"os"
)

func init() {
	go broadcastDataEvents(jsonEventChannel)
}

func broadcastDataEvents(in chan string) {
	for {
		os.Stdout.Write([]byte(<-in))
		os.Stdout.Write([]byte("\n"))
	}
}
