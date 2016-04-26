package errutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var testFactory = NewNamedFactory("testFactory", "test factory outside")

func Test_JSONStruct(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	factory := NewFactory("test error 1")
	jsondata, err := NewJSON(factory.New(New("test error 2")))
	require.NoError(err)
	require.NotNil(jsondata)

	data, err := json.Marshal(jsondata)
	require.NoError(err)
	require.Contains(string(data), `"test error 1"`)
	require.Contains(string(data), `"test error 2"`)
	require.Contains(string(data), `errutil/ErrorJSON_test.go:17`)
}

func Test_JSONStruct_outside(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	jsondata, err := NewJSON(testFactory.New(New("test error 2")))
	require.NoError(err)
	require.NotNil(jsondata)

	data, err := json.Marshal(jsondata)
	require.NoError(err)
	require.Contains(string(data), `"test factory outside"`)
	require.Contains(string(data), `"test error 2"`)
	require.Contains(string(data), `"testFactory"`)
}

func Test_JSONStruct_inherit(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	jsondata, err := NewJSON(NewErrors(
		New("test error 1"),
		New("test error 2"),
	))
	require.NoError(err)
	require.NotNil(jsondata)

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
}
