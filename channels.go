package report

// string-keyed map of unstructured data relevant to the event
type Data map[string]interface{}

//     info.go, action.go, timer.go
// --> rawEventChannel
// --> global.go <-- addGlobalChannel
// --> withGlobalsEventChannel
// --> json.go
// --> jsonEventChannel
// --> broadcast.go
var rawEventChannel = make(chan Data, 100)
var withGlobalsEventChannel = make(chan Data, 100)
var addGlobalChannel = make(chan Data)
var jsonEventChannel = make(chan string)
