package errutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_JSONErrorFormatter(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	formatter := NewJSONFormatter()

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	errtext, err := formatter.Format(errtest)
	require.NoError(err)
	require.Contains(errtext, `test error 1`)
	require.Contains(errtext, `test error 2`)
	require.Contains(errtext, `JSONFormatter_test.go:20`)
}
