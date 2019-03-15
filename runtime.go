package report

import (
	"runtime"
	"time"
)

// RuntimeEvery records runtime stats at the specified interval
//
//     log := report.New(report.StdOutJSON(), report.Data{"service": "myAppName"})
//     log.RuntimeEvery(time.Second*10)
//
func (l *Logger) RuntimeEvery(duration time.Duration) {
	l.wg.Add(1)
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

	statLoop:
		for {
			select {
			case <-ticker.C:
				l.Info("runtime", runtimeData())
			case <-l.stopC:
				break statLoop
			}
		}

		l.wg.Done()
	}()
}

func runtimeData() Data {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	return Data{
		"runtime.stack_mb":        float64(m.StackSys) / float64(1024*1024),
		"runtime.heap_mb":         float64(m.HeapAlloc) / float64(1024*1024),
		"runtime.goroutine_count": runtime.NumGoroutine(),
		"runtime.gc_pause_ms":     float64(m.PauseNs[(m.NumGC+255)%256]) / msPerNs,
	}
}
