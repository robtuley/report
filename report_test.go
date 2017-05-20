package report_test

import (
	"fmt"
	"os"

	"github.com/rainchasers/report"
)

func Example() {
	// setup logging output
	log := report.New(os.Stdout, report.Data{
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

	// Output:
	// {"event":"example.start","service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"event":"example.tick","sequence":0,"service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"event":"example.tick","sequence":1,"service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"event":"example.tick","sequence":2,"service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
	// {"event":"example.stop","service":"example","timestamp":"2017-05-20T21:00:24.2+01:00","type":"info"}
}
