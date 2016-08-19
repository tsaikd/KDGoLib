package logrusutil_test

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/tsaikd/KDGoLib/logrusutil"
)

func ExampleStackLogLevel() {
	// logger := logrus.New()
	logger := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrusutil.ConsoleLogFormatter{
			Flag: logrusutil.Llevel,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.InfoLevel,
	}
	defer func() {
		fmt.Println("after logger level", logger.Level)
	}()

	logger.Debugln("debug message")
	logger.Infoln("info message")

	defer logrusutil.StackLogLevel(logger, logrus.DebugLevel)()
	fmt.Println("stack logger level", logger.Level)

	logger.Debugln("debug message")
	logger.Infoln("info message")
	// Output:
	// [info] info message
	// stack logger level debug
	// [debug] debug message
	// [info] info message
	// after logger level info
}
