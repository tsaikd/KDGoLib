package errutil

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Trace(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err := New("new error")
	require.Error(err)

	buffer := &bytes.Buffer{}
	TraceWriter(err, buffer)
	require.Contains(buffer.String(), `util_test.go:14`)
	require.Contains(buffer.String(), `"error":"new error","errors":["new error"]`)
}
