package report

// SpanOption is a option to modify a span
type SpanOption func(s Span) Span

// TraceID option specifies the trace ID for a span
func TraceID(id string) SpanOption {
	return func(s Span) Span {
		s.traceID = id
		return s
	}
}

// ParentSpanID option specifies the parent span ID for a span
func ParentSpanID(id string) SpanOption {
	return func(s Span) Span {
		s.parentID = id
		return s
	}
}
