package report

// string-keyed map of unstructured data relevant to the event
type Data map[string]interface{}

//     info.go, action.go, timer.go
// --> rawEventChannel
// --> global.go <-- addGlobalChannel
//               --> drainChannel
// --> withGlobalsEventChannel
// --> json.go --> drainChannel
// --> jsonEventChannel
// --> broadcast.go --> drainChannel
var rawEventChannel = make(chan Data, 100)
var withGlobalsEventChannel = make(chan Data, 100)
var addGlobalChannel = make(chan Data)
var jsonEventChannel = make(chan string)
var drainChannel = make(chan bool)

// waits for events to drain down before exiting
func Drain() {
	close(rawEventChannel)
	<-drainChannel

	close(withGlobalsEventChannel)
	<-drainChannel

	close(jsonEventChannel)
	<-drainChannel
}
