package report

import (
	"encoding/json"
	"io"
	"os"

	"github.com/honeycombio/libhoney-go"
)

// Exporter exports events to an external service
type Exporter interface {
	Send(d Data) error
	Close()
}

// JSON writes JSON formatted logs
func JSON(w io.Writer) Exporter {
	return jw{
		encoder: json.NewEncoder(w),
	}
}

// StdOutJSON writes logs to StdOut as JSON
func StdOutJSON() Exporter {
	return JSON(os.Stdout)
}

// Honeycomb sends log events to HoneyComb
func Honeycomb(key string, dataset string) Exporter {
	libhoney.Init(libhoney.Config{
		WriteKey: key,
		Dataset:  dataset,
	})

	return hw{}
}

// json writer
type jw struct {
	encoder *json.Encoder
}

func (w jw) Send(d Data) error {
	return w.encoder.Encode(d)
}

func (w jw) Close() {}

// honeycomb writer
type hw struct{}

func (w hw) Send(d Data) error {
	ev := libhoney.NewEvent()
	ev.Add(d)
	return ev.Send()
}

func (w hw) Close() {
	libhoney.Close()
}
