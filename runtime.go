package report

import (
	"runtime"
	"time"
)

// RuntimeStatEvery records runtime stats at the specified interval
//
//     log := report.New(report.StdOutJSON(), report.Data{"service": "myAppName"})
//     log.RuntimeStatEvery("runtime", time.Second*10)
//
func (l *Logger) RuntimeStatEvery(event string, duration time.Duration) {
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

	statLoop:
		for {
			select {
			case <-ticker.C:
				l.Info(event, runtimeData())
			case <-l.stopC:
				break statLoop
			}
		}
	}()
}

func runtimeData() Data {
	data := Data{}

	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	data["stack.mb"] = float64(m.StackSys) / float64(1024*1024)
	data["heap.mb"] = float64(m.HeapAlloc) / float64(1024*1024)
	data["goroutine.count"] = runtime.NumGoroutine()
	data["gc.pause.ns"] = m.PauseNs[(m.NumGC+255)%256]

	return data
}
