package report

import (
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	Global(Data{"application": "myAppName"})
	for i := 0; i < b.N; i++ {
		Info("event.name", Data{"a": "aString", "z": 12})
	}
}
