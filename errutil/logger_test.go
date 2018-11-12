package errutil

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger_Print(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	logger := defaultLogger
	defer func() {
		defaultLogger = logger
	}()

	buffer := &bytes.Buffer{}
	SetDefaultLogger(loggerImpl{LoggerOptions{
		DefaultOutput:  buffer,
		TraceFormatter: NewJSONFormatter(),
	}})

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	Logger().Print(errtest.Error())

	errtext := buffer.String()
	require.Contains(errtext, `test error 1`)
	require.Contains(errtext, `test error 2`)
	require.Contains(errtext, `logger_test.go:33 `)
	require.True(strings.HasSuffix(errtext, "\n"))
}

func TestLogger_Error(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	logger := defaultLogger
	defer func() {
		defaultLogger = logger
	}()

	buffer := &bytes.Buffer{}
	SetDefaultLogger(loggerImpl{LoggerOptions{
		ErrorOutput:    buffer,
		TraceFormatter: NewJSONFormatter(),
	}})

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	Logger().Error(errtest.Error())

	errtext := buffer.String()
	require.Contains(errtext, `test error 1`)
	require.Contains(errtext, `test error 2`)
	require.Contains(errtext, `logger_test.go:63 `)
	require.True(strings.HasSuffix(errtext, "\n"))
}

func TestLogger_Trace(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	logger := defaultLogger
	defer func() {
		defaultLogger = logger
	}()

	buffer := &bytes.Buffer{}
	SetDefaultLogger(loggerImpl{LoggerOptions{
		ErrorOutput:    buffer,
		TraceFormatter: NewJSONFormatter(),
	}})

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	Logger().Trace(errtest)

	errtext := buffer.String()
	require.Contains(errtext, `"test error 1"`)
	require.Contains(errtext, `"test error 2"`)
	require.Contains(errtext, `logger_test.go:91`)
	require.True(strings.HasSuffix(errtext, "\n"))
}
