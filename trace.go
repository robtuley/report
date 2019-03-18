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

// Trace initialises a trace within the returned context
func (l *Logger) Trace(ctx context.Context) (traceCtx context.Context, traceID string) {
	s := fromContext(ctx)
	if s.TraceID == "" {
		s.TraceID = createULID()
	}
	return context.WithValue(ctx, contextKey, s), s.TraceID
}

// ContinueTrace injects trace ID into returned context
func (l *Logger) ContinueTrace(ctx context.Context, traceID string) context.Context {
	s := fromContext(ctx)
	s.TraceID = traceID
	return context.WithValue(ctx, contextKey, s)
}

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

	// if no required overall trace ID, generate this
	// but best results if there is either a root span, or
	// an injected trace/causation reference.
	if s.TraceID == "" {
		s.TraceID = createULID()
	}

	return context.WithValue(ctx, contextKey, s)
}

const msPerNs = float64(time.Millisecond) / float64(time.Nanosecond)

// EndSpan writes the data from a completed trace span
func (l *Logger) EndSpan(ctx context.Context, err error, payload Data) <-chan int {
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

	// if span has resulted in an unhandled error, flag event as actionable
	cmd := span
	if err != nil {
		payload["error"] = err.Error()
		cmd = action
	}

	// dispatch log task
	ack := make(chan int)
	l.taskC <- task{
		command: cmd,
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
