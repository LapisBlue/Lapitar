package util

import (
	"time"
)

type StopWatch struct {
	start   time.Time
	elapsed time.Duration
	marks   []time.Time
	running bool
}

func (watch *StopWatch) Start() *StopWatch {
	if !watch.running {
		watch.running = true
		watch.start = time.Now()
	}

	return watch
}

func (watch *StopWatch) Mark() *StopWatch {
	if watch.running {
		watch.marks = append(watch.marks, time.Now())
	}

	return watch
}

func (watch *StopWatch) Stop() *StopWatch {
	if watch.running {
		watch.marks = watch.marks[:0]
		watch.elapsed = watch.Elapsed()
		watch.running = false
	}

	return watch
}

func (watch *StopWatch) IsRunning() bool {
	return watch.running
}

func (watch *StopWatch) Elapsed() time.Duration {
	last := len(watch.marks)
	if last > 0 {
		result := time.Now().Sub(watch.marks[last-1])
		watch.marks = watch.marks[:last-1]
		return result
	}

	if watch.running {
		return time.Now().Sub(watch.start)
	} else {
		return watch.elapsed
	}
}

func (watch *StopWatch) String() string {
	return "(" + watch.Elapsed().String() + ")"
}

func StartedWatch() (watch *StopWatch) {
	watch = StoppedWatch()
	watch.Start()
	return
}

func StoppedWatch() (watch *StopWatch) {
	watch = new(StopWatch)
	return
}

var global *StopWatch

func GlobalWatch() *StopWatch {
	if global == nil {
		global = StoppedWatch()
	}

	return global
}
