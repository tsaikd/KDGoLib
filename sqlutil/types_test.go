package sqlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLValueStringSlice(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	testStringSlice := SQLStringSlice{"abc", "def"}
	if sqlvalue, err := SQLValueStringSlice(testStringSlice); assert.NoError(err) {
		require.Equal(`{"abc","def"}`, sqlvalue)
	}

	if sqlvalue, err := SQLValueStringSlice(&testStringSlice); assert.NoError(err) {
		require.Equal(`{"abc","def"}`, sqlvalue)
	}

	if sqlvalue, err := SQLValueStringSlice(nil); assert.NoError(err) {
		require.Nil(sqlvalue)
	}
}
