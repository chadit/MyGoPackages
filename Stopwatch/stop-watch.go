package stopwatch

import (
	"time"
)

// StopWatch -
type StopWatch struct {
	start, stop time.Time
}

// NewStopWatch - returns a new stopwatch
func NewStopWatch() StopWatch {
	return StopWatch{}
}

// Milliseconds - return the difference
func (stopWatch *StopWatch) Milliseconds() int {
	return int(stopWatch.stop.Sub(stopWatch.start) / time.Millisecond)
}

// Start - starts the stopwatch
func (stopWatch *StopWatch) Start() {
	stopWatch.start = time.Now().UTC()
}

// Stop - stops the stopwatch
func (stopWatch *StopWatch) Stop() {
	stopWatch.stop = time.Now().UTC()
}
