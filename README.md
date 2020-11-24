# Logging Utility for Go

A simple log & trace utility for Golang services.

[![Go Report Card](https://goreportcard.com/badge/github.com/robtuley/report)](https://goreportcard.com/report/github.com/robtuley/report)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/robtuley/report)](https://pkg.go.dev/github.com/robtuley/report)

- `Info` & `Action` logging methods to record information or actionable errors.
- `StartSpan` & `EndSpan` methods to record trace information.
- Runtime Golang stats are gathered via `RuntimeEvery`.
- Flexible writer interface provided, with StdOut JSON & Honeycomb exporters included.
- Log metrics aggregated and exposed via `Count` for log interface tests.
