package errorJson

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/errutil"
)

func Test_String(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	errjson, err := errutil.NewJSON(errutil.NewErrors(
		errors.New("test error 1"),
		errors.New("test error 2"),
	))
	require.NoError(err)
	reserr := &responseError{
		Status:    404,
		ErrorJSON: errjson,
	}
	data, err := json.Marshal(reserr)
	require.NoError(err)
	require.Contains(string(data), `404`)
	require.Contains(string(data), `"test error 1"`)
	require.Contains(string(data), `"test error 2"`)
	require.Contains(string(data), `errorJson_test.go:19`)

	store := map[string]interface{}{}
	err = json.Unmarshal(data, &store)
	require.NoError(err)
	require.Equal(float64(404), store["status"])
	require.Equal("test error 1", store["error"])
	require.Contains(store["errorpath"], `errorJson_test.go:19`)
}
