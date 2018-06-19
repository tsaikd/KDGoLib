package logrusutil_test

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/logrusutil"
)

func Test_ConsoleLogFormatter(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	outputBuffer := &bytes.Buffer{}
	logger := &logrus.Logger{
		Out:       outputBuffer,
		Formatter: &logrusutil.ConsoleLogFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	logger.Println("TEST OUTPUT")
	require.Contains(outputBuffer.String(), "ConsoleLogFormatter_test.go:24")
	require.Contains(outputBuffer.String(), "TEST OUTPUT")
}
