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

func TestSQLValueStringSliceJSON(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	value := SQLStringSliceJSON{}
	err := value.Scan([]byte(`["a","b","c"]`))
	require.NoError(err)
	require.Len(value, 3)
	require.Equal("a", value[0])
	require.Equal("b", value[1])
	require.Equal("c", value[2])

	sqlv, err := value.Value()
	require.NoError(err)
	require.Equal([]byte(`["a","b","c"]`), sqlv)
}

func TestSQLJsonMap(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	value := SQLJsonMap{}
	err := value.Scan([]byte(`{"k1":"str","k2":123}`))
	require.NoError(err)
	require.Equal("str", value["k1"])
	require.EqualValues(123, value["k2"])

	sqlv, err := value.Value()
	require.NoError(err)
	require.Equal([]byte(`{"k1":"str","k2":123}`), sqlv)
}
