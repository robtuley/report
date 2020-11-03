package report

import (
	"time"
)

// Span represents a trace span
type Span struct {
	traceID     string
	parentID    string
	spanID      string
	timestamp   time.Time
	duration    time.Duration
	isEnded     bool
	event       string
	err         error
	data        Data
	linkedSpans []Span
}

// StartSpan creates a span
func StartSpan(event string, opts ...SpanOption) Span {
	span := Span{
		traceID:   createULID(),
		spanID:    createULID(),
		timestamp: time.Now(),
		event:     event,
		data:      Data(make(map[string]interface{})),
	}
	for _, opt := range opts {
		span = opt(span)
	}
	return span
}

// TraceID provides the span trace ID
func (s Span) TraceID() string {
	return s.traceID
}

// SpanID provides the span ID
func (s Span) SpanID() string {
	return s.spanID
}

// Field associates key and value with the span
func (s Span) Field(k string, v interface{}) Span {
	s.data[k] = v
	return s
}

// End finishes the span
func (s Span) End(errors ...error) Span {

errLoop:
	for _, e := range errors {
		if e != nil {
			// break after finding the first non-nil err
			s.err = e
			break errLoop
		}
	}

	s.duration = time.Now().Sub(s.timestamp)
	s.isEnded = true
	return s
}

// Child adds a child span
func (s Span) Child(child Span) Span {
	child.traceID = s.traceID
	child.parentID = s.spanID
	s.linkedSpans = append(s.linkedSpans, child)
	return s
}

// FollowedBy adds a followed by span
func (s Span) FollowedBy(next Span) Span {
	next.traceID = s.traceID
	next.parentID = s.spanID
	next.linkedSpans = append(next.linkedSpans, s)
	return next
}

// Err provides error if there was one within span
func (s Span) Err() error {
	for _, span := range s.linkedSpans {
		if err := span.Err(); err != nil {
			return err
		}
	}
	return s.err
}
