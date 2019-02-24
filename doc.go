// Package report provides a simple log and tracing interface.
//
// Info & Action logging methods to record information or actionable errors.
// StartSpan & EndSpan methods to record trace information.
// Runtime Golang stats are gathered via RuntimeEvery.
// Flexible writer interface provided, with StdOut JSON & Honeycomb exporters included.
// Log metrics aggregated and exposed via Count for log interface tests.
// Default package level logger to simplify your call chain.
//
// See Also
//
// https://opentracing.io/docs/overview/spans/
// https://docs.honeycomb.io/working-with-data/tracing/
// https://docs.honeycomb.io/working-with-data/tracing/send-trace-data/
package report
