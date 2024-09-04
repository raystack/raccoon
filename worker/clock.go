package worker

import "time"

type Clock interface {
	Now() time.Time
}

type defaultClock struct{}

func (defaultClock) Now() time.Time {
	return time.Now()
}

var DefaultClock = defaultClock{}
