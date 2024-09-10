package clock

import (
	"time"
)

// Clock represents a time source
// It can be used as a way to abstract away the implicit
// dependency on time.Now() from code.
type Clock interface {
	Now() time.Time
}

type defaultClock struct{}

func (defaultClock) Now() time.Time {
	return time.Now()
}

// Default clock. Uses time.Now() internally.
var Default = defaultClock{}
