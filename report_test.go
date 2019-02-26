package report_test

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/rainchasers/report"
)

func Example() {
	// setup logging output
	log := report.New("example")
	log.Baggage("timestamp", "2017-05-20T21:00:24.2+01:00") // to make output deterministic
	log.Export(report.StdOutJSON())
	defer log.Close()

	// ticker daemon execution
	log.Info("example.start", report.Data{})
	seq := 0
	for i := 0; i < 3; i++ {
		log.Info("example.tick", report.Data{"sequence": seq})
		seq = seq + 1
	}
	log.Info("example.stop", report.Data{})

	// validate
	if log.Count("example.tick") != 3 {
		// your own log validation...
		fmt.Print("Ooops! example.tick should be 3")
	}
	if err := log.Err(); err != nil {
		// your own log validation...
		fmt.Print(err)
	}

	// Output:
	// {"name":"example.start","service_name":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"example.tick","sequence":0,"service_name":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"example.tick","sequence":1,"service_name":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"example.tick","sequence":2,"service_name":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"example.stop","service_name":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
}

func ExampleLogger_Info() {
	// setup logging output
	log := report.New("example")
	log.Baggage("timestamp", "2017-05-20T21:00:24.2+01:00") // to make output deterministic
	log.Export(report.StdOutJSON())
	defer log.Close()

	// normal usage is simple to call log.Info
	log.Info("http.response", report.Data{"status": 200, "request": "/page1"})
	log.Info("http.response", report.Data{"status": 200, "request": "/page2"})

	// if you want to block until the logline is written, consume from the returned channel
	<-log.Info("http.response", report.Data{"status": 404, "request": "/nopage"})

	if err := log.Err(); err != nil {
		fmt.Print(err)
	}

	// Output:
	// {"name":"http.response","request":"/page1","service_name":"example","status":200,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"http.response","request":"/page2","service_name":"example","status":200,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"http.response","request":"/nopage","service_name":"example","status":404,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
}

func ExampleLogger_Action() {
	// setup logging output
	log := report.New("example")
	log.Baggage("timestamp", "2017-05-20T21:00:24.2+01:00") // to make output deterministic
	log.Export(report.StdOutJSON())
	defer log.Close()

	// for example we get an error...
	err := errors.New("Failed to parse JSON")

	// normal usage is simple to call log.Action
	log.Action("json.unparseable", report.Data{"error": err.Error()})

	// if you want to block until the logline is written, consume from the returned channel
	// (useful if you intend to shutdown as a result of this error)
	<-log.Action("json.unparseable", report.Data{"error": err.Error()})

	// LastError can be used to validate if an actionable event was logged
	if err := log.Err(); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// {"error":"Failed to parse JSON","name":"json.unparseable","service_name":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"action"}
	// {"error":"Failed to parse JSON","name":"json.unparseable","service_name":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"action"}
	// Actionable event: json.unparseable
}

func ExampleLogger_Count() {
	log := report.New("example")
	defer log.Close()

	log.Info("http.response.200", report.Data{})
	log.Info("http.response.404", report.Data{})
	log.Info("http.response.200", report.Data{})

	fmt.Printf("404 response count is %d", log.Count("http.response.404"))

	if err := log.Err(); err != nil {
		fmt.Print(err)
	}

	// Output:
	// 404 response count is 1
}

func ExampleLogger_Err() {
	log := report.New("example")
	log.Export(report.StdOutJSON())
	defer log.Close()

	// this log line has errored, but to prevent error handling clutter
	// the library interface requires you need to check error states
	// separately using Logger.LastError()
	<-log.Info("encoding.fail", report.Data{
		"unencodeable": make(chan int),
	})

	// check whether there has been any logging errors
	if err := log.Err(); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// Error sending encoding.fail: json: unsupported type: chan int
}

func ExampleLogger_Export() {
	// want to write to 2 logfiles, represented here as 2 buffers
	b1 := &bytes.Buffer{}
	b2 := &bytes.Buffer{}

	// setup logging output
	log := report.New("example")
	log.Baggage("timestamp", "2017-05-20T21:00:24.2+01:00") // to make output deterministic
	defer log.Close()

	// configure 2 writers
	log.Export(report.JSON(b1))
	log.Export(report.JSON(b2))

	// log something
	<-log.Info("http.response", report.Data{"status": 404, "request": "/nopage"})

	// output the 2 log files, note these have been written in parallel which
	// is why they are kept separate in this example until the end
	fmt.Print(b1)
	fmt.Print(b2)

	if err := log.Err(); err != nil {
		fmt.Print(err)
	}

	// Output:
	// {"name":"http.response","request":"/nopage","service_name":"example","status":404,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"http.response","request":"/nopage","service_name":"example","status":404,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
}
