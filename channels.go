package report

// string-keyed map of unstructured data relevant to the event
type Data map[string]interface{}

//     info.go, action.go, timer.go
// --> channel.RawEvents
// --> global.go <-- channel.AddGlobal
//               --> channel.Drain
// --> channel.WithGlobals
// --> json.go --> channel.Drain
// --> channel.JsonEncoded
// --> broadcast.go --> channel.Drain
var channel struct {
	RawEvents   chan Data
	WithGlobals chan Data
	AddGlobal   chan Data
	JsonEncoded chan string
	Drain       chan bool
}

func init() {
	channel.RawEvents = make(chan Data, 50)
	channel.WithGlobals = make(chan Data, 50)
	channel.AddGlobal = make(chan Data)
	channel.JsonEncoded = make(chan string, 50)
	channel.Drain = make(chan bool)
}

// waits for events to drain down before exiting
func Drain() {
	close(channel.RawEvents)
	<-channel.Drain

	close(channel.WithGlobals)
	<-channel.Drain

	close(channel.JsonEncoded)
	<-channel.Drain
}
