package logger

import (
	"io"
	"os"
	"raccoon/config"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func Setup() {
	if logger != nil {
		return
	}
	logLevel, err := log.ParseLevel(config.LogLevel())
	if err != nil {
		log.Panic(err)
	}
	logger = &log.Logger{
		Out: os.Stdout,
		Formatter: &log.TextFormatter{
			FullTimestamp: true,
		},
		Hooks: make(log.LevelHooks),
		Level: logLevel,
	}

	return
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
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
