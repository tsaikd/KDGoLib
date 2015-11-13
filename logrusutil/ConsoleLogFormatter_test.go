package logrusutil_test

import (
	"bytes"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tsaikd/KDGoLib/logrusutil"
)

func Test_ConsoleLogFormatter(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	outputBuffer := &bytes.Buffer{}
	logger := &logrus.Logger{
		Out:       outputBuffer,
		Formatter: &logrusutil.ConsoleLogFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	logger.Println("TEST OUTPUT")
	assert.Contains(outputBuffer.String(), "ConsoleLogFormatter_test.go:")
	assert.Contains(outputBuffer.String(), "TEST OUTPUT")
}
