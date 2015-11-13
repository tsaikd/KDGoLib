package logrusutil

import (
	"os"

	"github.com/Sirupsen/logrus"
)

func NewConsoleLogger() *logrus.Logger {
	return &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &ConsoleLogFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
}
