package report

import (
	"time"
)

const msPerNs = float64(time.Millisecond) / float64(time.Nanosecond)

// Trace writes the data from a completed trace span
func (l *Logger) Trace(s Span) <-chan int {
	for i := len(s.linkedSpans) - 1; i >= 0; i-- {
		l.sendSpan(s.linkedSpans[i])
	}
	// ignoring the logic error here where if the root
	// span is not ended, the returned channel won't block
	// on children writing
	return l.sendSpan(s)
}

func (l *Logger) sendSpan(s Span) <-chan int {
	if !s.isEnded {
		ch := make(chan int)
		close(ch)
		return ch
	}

	payload := s.data

	// add completed span data
	payload["duration_ms"] = float64(time.Now().Sub(s.timestamp).Nanoseconds()) / msPerNs
	payload["trace.span_id"] = s.spanID
	payload["trace.parent_id"] = s.parentID
	payload["trace.trace_id"] = s.traceID
	payload["timestamp"] = s.timestamp.Format(time.RFC3339Nano)

	// if span has resulted in an unhandled error, flag event as actionable
	cmd := spanCmd
	if s.err != nil {
		payload["error"] = s.err.Error()
		cmd = action
	}

	// dispatch log task
	ack := make(chan int)
	l.taskC <- task{
		command: cmd,
		event:   s.event,
		data:    payload,
		ackC:    ack,
	}

	return ack
}
