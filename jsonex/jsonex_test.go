package jsonex_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/jsonex"
)

func TestIsEmpty(t *testing.T) {
	type myStruct struct {
		String string
	}

	require := require.New(t)
	require.NotNil(require)

	require.False(jsonex.IsEmpty(true))
	require.False(jsonex.IsEmpty(1))
	require.False(jsonex.IsEmpty("text"))
	require.False(jsonex.IsEmpty(myStruct{String: "text"}))

	require.True(jsonex.IsEmpty(false))
	require.True(jsonex.IsEmpty(0))
	require.True(jsonex.IsEmpty(""))
	require.True(jsonex.IsEmpty([]interface{}{}))
	require.True(jsonex.IsEmpty(marshalStruct{"text"}))
	require.True(jsonex.IsEmpty(myStruct{}))

	require.True(jsonex.IsEmpty((*reflect.Value)(nil)))
	require.True(jsonex.IsEmpty(reflect.ValueOf(false)))
	if v := reflect.ValueOf(false); v.IsValid() {
		require.True(jsonex.IsEmpty(&v))
	}
}
