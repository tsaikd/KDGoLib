package logutil

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/inject"
	"github.com/tsaikd/KDGoLib/logrusutil"
)

var (
	timestampFormat = "2006/01/02 15:04:05"
	DefaultLogger   = &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrusutil.ConsoleLogFormatter{
			TimestampFormat: timestampFormat,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.InfoLevel,
	}
)

type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

type LevelLogger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Debugln(...interface{})

	Info(...interface{})
	Infof(string, ...interface{})
	Infoln(...interface{})

	Warn(...interface{})
	Warnf(string, ...interface{})
	Warnln(...interface{})

	Error(...interface{})
	Errorf(string, ...interface{})
	Errorln(...interface{})
}

// GetStdLogger get StdLogger instance from injector
func GetStdLogger(inj inject.Injector) StdLogger {
	if loggerType := inj.Get(inject.InterfaceOf((*StdLogger)(nil))); loggerType.IsValid() {
		logger := loggerType.Interface().(StdLogger)
		return logger
	}
	return DefaultLogger
}

// GetLevelLogger get LevelLogger instance from injector
func GetLevelLogger(inj inject.Injector) LevelLogger {
	if loggerType := inj.Get(inject.InterfaceOf((*LevelLogger)(nil))); loggerType.IsValid() {
		logger := loggerType.Interface().(LevelLogger)
		return logger
	}
	return DefaultLogger
}
