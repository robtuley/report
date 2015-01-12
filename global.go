package report

func Global(payload Data) {
	addGlobalChannel <- payload
}

func init() {
	go addGlobalsToEventStream(rawEventChannel, withGlobalsEventChannel, addGlobalChannel)
}

func addGlobalsToEventStream(in chan Data, out chan Data, add chan Data) {
	globals := Data{}

	for {
		select {
		case evt := <-in:
			for k, v := range globals {
				evt[k] = v
			}
			out <- evt
		case g := <-add:
			for k, v := range g {
				globals[k] = v
			}
		}
	}
}
