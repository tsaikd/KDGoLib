package errutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_JSONStruct(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	jsondata := NewJSON(NewErrors(
		New("test error 1"),
		New("test error 2"),
	))

	data, err := json.Marshal(jsondata)
	require.NoError(err)
	require.Contains(string(data), `"test error 1"`)
	require.Contains(string(data), `"test error 2"`)
	require.Contains(string(data), `ErrorJSON_test.go:15`)
}

func Test_JSONStruct_inherit(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	jsondata := NewJSON(NewErrors(
		New("test error 1"),
		New("test error 2"),
	))

	teststruct := struct {
		*ErrorJSON
		Field string `json:"field"`
	}{
		ErrorJSON: jsondata,
		Field:     "test field",
	}

	data, err := json.Marshal(teststruct)
	require.NoError(err)
	require.Contains(string(data), `"test field"`)
	require.Contains(string(data), `"test error 1"`)
	require.Contains(string(data), `"test error 2"`)
	require.Contains(string(data), `ErrorJSON_test.go:31`)
}
