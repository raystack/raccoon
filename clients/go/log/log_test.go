package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type customLogger struct{}

func (c customLogger) Infof(msg string, keysAndValues ...interface{}) {}

func (c customLogger) Errorf(msg string, keysAndValues ...interface{}) {}

func TestRestLogger(t *testing.T) {
	assert := assert.New(t)

	var l interface{} = &customLogger{}
	l, ok := l.(Logger)

	assert.True(ok)
	assert.NotNil(l)
}
