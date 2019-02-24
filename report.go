package report

import (
	"errors"
	"sync"
	"time"

	"github.com/honeycombio/libhoney-go"
)

// Logger is the central logging agent on which to register events
type Logger struct {
	writer         func(d Data) error
	taskC          chan task
	stopC          chan struct{}
	global         Data
	count          map[string]int
	lastError      error
	lastErrorMutex sync.Mutex
}

// Data is a string-keyed map of unstructured data relevant to the event
type Data map[string]interface{}

type command int

const (
	info command = iota
	action
	span
	count
)

type task struct {
	command command
	event   string
	data    Data
	ackC    chan<- int
}

// New creates an instance of a logging agent
//
//     logger := report.New(report.JSON(), report.Data{"service": "myAppName"})
//     defer logger.Stop()
//
func New(w Writer, global Data) *Logger {
	logger := Logger{
		writer: w,
		taskC:  make(chan task, 1),
		stopC:  make(chan struct{}),
		global: global,
		count:  make(map[string]int),
	}
	go logger.run()
	return &logger
}

// Info logs event that provide telemetry measures or context to any events requiring action.
func (l *Logger) Info(event string, payload Data) <-chan int {
	ack := make(chan int)
	l.taskC <- task{
		command: info,
		event:   event,
		data:    payload,
		ackC:    ack,
	}
	return ack
}

// Action events that need intervention or resolving.
func (l *Logger) Action(event string, payload Data) <-chan int {
	ack := make(chan int)
	l.taskC <- task{
		command: action,
		event:   event,
		data:    payload,
		ackC:    ack,
	}

	return ack
}

// Count returns the number of log events of a particular type since startup
func (l *Logger) Count(event string) int {
	ack := make(chan int)
	l.taskC <- task{
		command: count,
		event:   event,
		data:    Data{},
		ackC:    ack,
	}

	return <-ack
}

// LastError returns the last Actionable log event or encoding error if either occurred
func (l *Logger) LastError() error {
	l.lastErrorMutex.Lock()
	defer l.lastErrorMutex.Unlock()

	return l.lastError
}

// Stop shuts down the logging agent, further logging will result in a panic
//
//     log := report.New(report.JSON(), report.Data{"service": "myAppName"})
//     defer log.Stop()
//
func (l *Logger) Stop() {
	close(l.taskC)
	close(l.stopC)
	// we should call libhoney.Close() here but if not Inited this
	// causes a panic. So we call Flush() to ensure any pending events
	// have been sent.
	libhoney.Flush()
}

func (l *Logger) run() {

toNewTask:
	for t := range l.taskC {
		if t.command == count {
			n, exists := l.count[t.event]
			if exists {
				t.ackC <- n
			} else {
				t.ackC <- 0
			}
			close(t.ackC)
			continue toNewTask
		}

		n, exists := l.count[t.event]
		if exists {
			l.count[t.event] = n + 1
		} else {
			l.count[t.event] = 1
		}

		t.data["name"] = t.event
		if _, exists := t.data["timestamp"]; !exists {
			t.data["timestamp"] = time.Now().Format(time.RFC3339Nano)
		}
		for k, v := range l.global {
			t.data[k] = v
		}
		switch t.command {
		case info:
			t.data["type"] = "info"
		case action:
			t.data["type"] = "action"
			l.lastErrorMutex.Lock()
			l.lastError = errors.New("Actionable event: " + t.event)
			l.lastErrorMutex.Unlock()
		case span:
			t.data["type"] = "span"
		}

		if err := l.writer(t.data); err != nil {
			l.lastErrorMutex.Lock()
			l.lastError = err
			l.lastErrorMutex.Unlock()
		}
		close(t.ackC)
	}
}
