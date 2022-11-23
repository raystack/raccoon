package log

import (
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	Info LogLevel = iota
	Error
)

var (
	level LogLevel = Info

	// defaultLogger is the logger provided with client
	defaultLogger Logger = &consoleLogger{
		info: log.New(os.Stderr, "[INFO] ", log.LstdFlags),
		err:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}

	logInit sync.Once
)

type Logger interface {
	Infof(msg string, keysAndValues ...interface{})
	Errorf(msg string, keysAndValues ...interface{})
}

func Default() Logger {
	return defaultLogger
}

func SetLogger(logger Logger) {
	logInit.Do(func() {
		if logger == nil {
			defaultLogger = &nopLogger{}
			return
		}

		defaultLogger = logger
	})
}

func Infof(msg string, keysAndValues ...interface{}) {
	if defaultLogger == nil {
		return
	}

	defaultLogger.Infof(msg, keysAndValues...)
}

func Errorf(msg string, keysAndValues ...interface{}) {
	if defaultLogger == nil {
		return
	}

	defaultLogger.Errorf(msg, keysAndValues...)
}

// consoleLogger writes to the standart output.
type consoleLogger struct {
	info, err *log.Logger
}

func (c *consoleLogger) Infof(msg string, keysAndValues ...interface{}) {
	if level <= Info {
		c.info.Printf(msg, keysAndValues...)
	}
}

func (c *consoleLogger) Errorf(msg string, keysAndValues ...interface{}) {
	if level <= Error {
		c.err.Printf(msg, keysAndValues...)
	}
}

// nopLogger is the empty logger.
type nopLogger struct{}

func (*nopLogger) Infof(msg string, keysAndValues ...interface{}) {}

func (*nopLogger) Errorf(msg string, keysAndValues ...interface{}) {}
