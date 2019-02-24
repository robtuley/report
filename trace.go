package report

import (
	"context"
	"time"
)

type spanContext struct {
	TraceID   string
	ParentID  string
	SpanID    string
	Timestamp time.Time
	Name      string
}

type contextField int

const (
	contextKey contextField = iota
)

// TODO: func to inject traceID & parentID from external input (e.g. headers or msg)

// StartSpan marks the start of a trace span
func (l *Logger) StartSpan(ctx context.Context, event string) context.Context {
	s := fromContext(ctx)

	// if a SpanID is present there is an open active span,
	// so this SpanID needs to be reflected as the new ParentID
	if s.SpanID != "" {
		s.ParentID = s.SpanID
	}

	// init id, ts & name data fields for span
	s.SpanID = createULID()
	s.Timestamp = time.Now()
	s.Name = event

	// if no parent span, bubble up current span
	if s.ParentID == "" {
		s.ParentID = s.SpanID
		if s.TraceID == "" {
			s.TraceID = s.ParentID
		}
	}

	return context.WithValue(ctx, contextKey, s)
}

const msPerNs = float64(time.Millisecond) / float64(time.Nanosecond)

// EndSpan writes the data from a completed trace span
func (l *Logger) EndSpan(ctx context.Context, payload Data) <-chan int {
	s := fromContext(ctx)

	// check a span is open, if not simply log as error
	if s.SpanID == "" {
		payload["error"] = "End span called on an unopen span"
		return l.Action("trace.error", payload)
	}

	// write completed span data
	payload["duration_ms"] = float64(time.Now().Sub(s.Timestamp).Nanoseconds()) / msPerNs
	payload["trace.span_id"] = s.SpanID
	payload["trace.parent_id"] = s.ParentID
	payload["trace.trace_id"] = s.TraceID
	payload["timestamp"] = s.Timestamp.Format(time.RFC3339Nano)

	ack := make(chan int)
	l.taskC <- task{
		command: span,
		event:   s.Name,
		data:    payload,
		ackC:    ack,
	}

	return ack
}

func fromContext(ctx context.Context) spanContext {
	new := spanContext{}
	if ctx == nil {
		return new
	}
	if existing, ok := ctx.Value(contextKey).(spanContext); ok {
		return existing
	}
	return new
}
