package util

import (
	"time"
)

type StopWatch struct {
	start   time.Time
	elapsed time.Duration
	mark    *time.Time
	running bool
}

func (watch *StopWatch) Start() bool {
	if watch.running {
		return false
	}

	watch.running = true
	watch.start = time.Now()
	return true
}

func (watch *StopWatch) Mark() {
	now := time.Now()
	watch.mark = &now
}

func (watch *StopWatch) Stop() bool {
	if !watch.running {
		return false
	}

	watch.running = false
	watch.elapsed += time.Now().Sub(watch.start)
	return true
}

func (watch *StopWatch) Reset() *StopWatch {
	watch.running = false
	watch.elapsed = 0
	return watch
}

func (watch *StopWatch) IsRunning() bool {
	return watch.running
}

func (watch *StopWatch) Elapsed() time.Duration {
	if watch.mark != nil {
		result := time.Now().Sub(*watch.mark)
		watch.mark = nil
		return result
	}

	return watch.elapsed
}

func (watch *StopWatch) String() string {
	return "(" + watch.Elapsed().String() + ")"
}

func CreateStopWatch() (watch *StopWatch) {
	watch = new(StopWatch)
	return
}
