package worker

import "time"

type TimeSource interface {
	Now() time.Time
}

type defaultTimeSource struct{}

func (defaultTimeSource) Now() time.Time {
	return time.Now()
}

var DefaultTimeSource = defaultTimeSource{}
