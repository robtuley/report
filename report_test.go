package report_test

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/rainchasers/report"
)

func Example() {
	// setup logging output
	log := report.New(report.StdOutJSON(), report.Data{
		"service":   "example",
		"timestamp": "2017-05-20T21:00:24.2+01:00", // to make output deterministic
	})
	defer log.Stop()

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
	if err := log.LastError(); err != nil {
		// your own log validation...
		fmt.Print(err)
	}

	// Output:
	// {"name":"example.start","service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"example.tick","sequence":0,"service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"example.tick","sequence":1,"service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"example.tick","sequence":2,"service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"example.stop","service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
}

func ExampleLogger_Info() {
	// setup logging output, NOTE timestamp included only to make output deterministic
	log := report.New(report.StdOutJSON(), report.Data{"timestamp": "2017-05-20T21:00:24.2+01:00"})
	defer log.Stop()

	// normal usage is simple to call log.Info
	log.Info("http.response", report.Data{"status": 200, "request": "/page1"})
	log.Info("http.response", report.Data{"status": 200, "request": "/page2"})

	// if you want to block until the logline is written, consume from the returned channel
	<-log.Info("http.response", report.Data{"status": 404, "request": "/nopage"})

	// Output:
	// {"name":"http.response","request":"/page1","status":200,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"http.response","request":"/page2","status":200,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"http.response","request":"/nopage","status":404,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
}

func ExampleLogger_Action() {
	// setup logging output, NOTE timestamp included only to make output deterministic
	log := report.New(report.StdOutJSON(), report.Data{"timestamp": "2017-05-20T21:00:24.2+01:00"})
	defer log.Stop()

	// for example we get an error...
	err := errors.New("Failed to parse JSON")

	// normal usage is simple to call log.Action
	log.Action("json.unparseable", report.Data{"error": err.Error()})

	// if you want to block until the logline is written, consume from the returned channel
	// (useful if you intend to shutdown as a result of this error)
	<-log.Action("json.unparseable", report.Data{"error": err.Error()})

	// LastError can be used to validate if an actionable event was logged
	if err := log.LastError(); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// {"error":"Failed to parse JSON","name":"json.unparseable","timestamp":"2017-05-20T21:00:24.2+01:00","type":"action"}
	// {"error":"Failed to parse JSON","name":"json.unparseable","timestamp":"2017-05-20T21:00:24.2+01:00","type":"action"}
	// Actionable event: json.unparseable
}

func ExampleLogger_Count() {
	discard := func(d report.Data) error {
		return nil
	}

	log := report.New(discard, report.Data{})
	defer log.Stop()

	log.Info("http.response.200", report.Data{})
	log.Info("http.response.404", report.Data{})
	log.Info("http.response.200", report.Data{})

	fmt.Printf("404 response count is %d", log.Count("http.response.404"))

	// Output:
	// 404 response count is 1
}

func ExampleLogger_LastError() {
	// this is a writer that will report an error
	writeError := func(d report.Data) error {
		return errors.New("a send error")
	}

	log := report.New(writeError, report.Data{})
	defer log.Stop()

	// this log line has errored, but to prevent error handling clutter
	// the library interface requires you need to check error states
	// separately using Logger.LastError()
	<-log.Info("encoding.fail", report.Data{})

	// check whether there has been any logging errors
	if err := log.LastError(); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// a send error
}

func ExampleWriter_And() {
	// want to write to 2 logfiles, represented here as 2 buffers
	b1 := &bytes.Buffer{}
	b2 := &bytes.Buffer{}

	// use Writer.And to create a single writer from 2 writers
	w := report.JSON(b1).And(report.JSON(b2))

	// setup logging output, NOTE timestamp included only to make output deterministic
	log := report.New(w, report.Data{"timestamp": "2017-05-20T21:00:24.2+01:00"})
	defer log.Stop()

	// log something
	<-log.Info("http.response", report.Data{"status": 404, "request": "/nopage"})

	// output the 2 log files, note these have been written in parallel which
	// is why they are kept separate in this example until the end
	fmt.Print(b1)
	fmt.Print(b2)

	// Output:
	// {"name":"http.response","request":"/nopage","status":404,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"name":"http.response","request":"/nopage","status":404,"timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
}
