package errutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_JSONErrorFormatter(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	formatter := NewBufferErrorFormatter()
	require.NotNil(formatter)

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	err := formatter.Format(errtest)
	require.NoError(err)

	errtext := formatter.String()
	require.Contains(errtext, `"test error 1"`)
	require.Contains(errtext, `"test error 2"`)
	require.Contains(errtext, `ErrorFormatter_test.go:18`)
}
