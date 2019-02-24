package report

import (
	"encoding/json"
	"io"
	"os"

	"github.com/honeycombio/libhoney-go"
)

// Writer is a event writing function
type Writer func(d Data) error

// And adds another writer to execute in parallel
func (w Writer) And(next Writer) Writer {
	return func(d Data) error {
		ch := make(chan error)
		go func() {
			ch <- w(d)
		}()
		go func() {
			ch <- next(d)
		}()
		if err := <-ch; err != nil {
			<-ch
			return err
		}
		return <-ch
	}
}

// JSON writes JSON formatted logs
func JSON(w io.Writer) Writer {
	encoder := json.NewEncoder(w)
	return func(d Data) error {
		return encoder.Encode(d)
	}
}

// StdOutJSON writes logs to StdOut as JSON
func StdOutJSON() Writer {
	return JSON(os.Stdout)
}

// Honeycomb sends log events to HoneyComb
func Honeycomb(key string, dataset string) Writer {
	libhoney.Init(libhoney.Config{
		WriteKey: key,
		Dataset:  dataset,
	})

	return func(d Data) error {
		ev := libhoney.NewEvent()
		ev.Add(d)
		return ev.Send()
	}
}
