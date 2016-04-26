package errutil

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type invalidFormatter struct {
}

func (t *invalidFormatter) Format(errin error) (errtext string, err error) {
	return "", New("invalid formatter")
}

func Test_Trace(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	formatter := defaultFormatter
	traceOutput := defaultTraceOutput
	defer func() {
		defaultFormatter = formatter
		defaultTraceOutput = traceOutput
	}()

	buffer := &bytes.Buffer{}
	SetDefaultFormatter(NewJSONFormatter())
	SetDefaultTraceOutput(buffer)

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	Trace(errtest)

	errtext := buffer.String()
	require.Contains(errtext, `"test error 1"`)
	require.Contains(errtext, `"test error 2"`)
	require.Contains(errtext, `trace_test.go:35`)
}

func Test_TracePanic(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	formatter := &invalidFormatter{}
	require.NotNil(formatter)

	SetDefaultFormatter(formatter)
	SetDefaultTraceOutput(os.Stderr)

	errtest := NewErrors(
		New("test error 1"),
		New("test error 2"),
	)

	require.Panics(func() {
		Trace(errtest)
	})
}
