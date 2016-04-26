package errutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConsoleErrorFormatter(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	formatter := NewConsoleFormatter("")

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	errtext, err := formatter.Format(errtest)
	require.NoError(err)
	require.Equal(`test error`, errtext)
}

func Test_ConsoleErrorFormatter_seperator(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	formatter := NewConsoleFormatter("; ")

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	errtext, err := formatter.Format(errtest)
	require.NoError(err)
	require.Equal(`test error; test error 1; test error 2`, errtext)
}
