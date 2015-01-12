package report

import (
	"log"
)

func init() {
	go broadcastDataEvents(jsonEventChannel)
}

func broadcastDataEvents(in chan string) {
	for {
		log.Println(<-in)
	}
}
