package logger

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger
var defaultLevel = log.InfoLevel

func init() {
	if logger != nil {
		return
	}
	logger = &log.Logger{
		Out: os.Stdout,
		Formatter: &log.TextFormatter{
			FullTimestamp: true,
		},
		Hooks: make(log.LevelHooks),
		Level: defaultLevel,
	}

	return
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Set(log *log.Logger) {
	logger = log
}

func SetOutput(out io.Writer) {
	logger.SetOutput(out)
}

func SetLevel(level string) {
	if l, err := log.ParseLevel(level); err == nil {
		logger.SetLevel(l)
	} else {
		fmt.Printf("[logger] NoOps, Fail to parse log level: %v\n", err)
	}
}
