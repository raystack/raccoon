package clock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// Mock clock for testing
type Mock struct {
	mock.Mock
}

func (m *Mock) Now() time.Time {
	return m.Called().Get(0).(time.Time)
}
