package report

import (
	"runtime"
)

func getRuntimeData() Data {
	data := Data{}

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	data["heap"] = float64(memStats.HeapAlloc)
	data["goroutines"] = float64(runtime.NumGoroutine())

	return data
}
