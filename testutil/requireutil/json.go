package requireutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorUnsupportedJSONType1 = errutil.NewFactory("unsupported json type: %T")
)

// UnmarshalJSONMap unmarshal jsonobj to map[string]interface{}
func UnmarshalJSONMap(jsonobj interface{}) (jsonmap map[string]interface{}, err error) {
	jsonmap = map[string]interface{}{}
	switch jsonobj.(type) {
	case []byte:
		err = json.Unmarshal(jsonobj.([]byte), &jsonmap)
		return
	case *[]byte:
		return UnmarshalJSONMap(*jsonobj.(*[]byte))
	case json.Marshaler:
		var jsonbyte []byte
		jsonbyte, err = jsonobj.(json.Marshaler).MarshalJSON()
		if err != nil {
			return
		}
		return UnmarshalJSONMap(jsonbyte)
	}
	return nil, ErrorUnsupportedJSONType1.New(nil, jsonobj)
}

// RequireJSONFieldEqual field of jsonobj should equal to expected
func RequireJSONFieldEqual(t *testing.T, expected interface{}, jsonobj interface{}, field string, msgAndArgs ...interface{}) {
	require := require.New(t)
	require.NotNil(require)

	jsonmap, err := UnmarshalJSONMap(jsonobj)
	require.NoError(err, msgAndArgs...)

	require.Contains(jsonmap, field, msgAndArgs...)
	require.Equal(expected, jsonmap[field], msgAndArgs...)
}

// RequireJSONFieldNotEqual field of jsonobj should not equal to expected
func RequireJSONFieldNotEqual(t *testing.T, expected interface{}, jsonobj interface{}, field string, msgAndArgs ...interface{}) {
	require := require.New(t)
	require.NotNil(require)

	jsonmap, err := UnmarshalJSONMap(jsonobj)
	require.NoError(err, msgAndArgs...)

	require.Contains(jsonmap, field, msgAndArgs...)
	require.NotEqual(expected, jsonmap[field], msgAndArgs...)
}

// RequireJSONFieldEqualValues the value of field of jsonobj should equal to expected
func RequireJSONFieldEqualValues(t *testing.T, expected interface{}, jsonobj interface{}, field string, msgAndArgs ...interface{}) {
	require := require.New(t)
	require.NotNil(require)

	jsonmap, err := UnmarshalJSONMap(jsonobj)
	require.NoError(err, msgAndArgs...)

	require.Contains(jsonmap, field, msgAndArgs...)
	require.EqualValues(expected, jsonmap[field], msgAndArgs...)
}

// RequireJSONFieldExist field of jsonobj should be exist
func RequireJSONFieldExist(t *testing.T, jsonobj interface{}, field string, msgAndArgs ...interface{}) {
	require := require.New(t)
	require.NotNil(require)

	jsonmap, err := UnmarshalJSONMap(jsonobj)
	require.NoError(err, msgAndArgs...)

	require.Contains(jsonmap, field, msgAndArgs...)
}

// RequireJSONFieldNotExist field of jsonobj should not be exist
func RequireJSONFieldNotExist(t *testing.T, jsonobj interface{}, field string, msgAndArgs ...interface{}) {
	require := require.New(t)
	require.NotNil(require)

	jsonmap, err := UnmarshalJSONMap(jsonobj)
	require.NoError(err, msgAndArgs...)

	require.NotContains(jsonmap, field, msgAndArgs...)
}

// RequireJSONFieldNil field of jsonobj should be nil
func RequireJSONFieldNil(t *testing.T, jsonobj interface{}, field string, msgAndArgs ...interface{}) {
	require := require.New(t)
	require.NotNil(require)

	jsonmap, err := UnmarshalJSONMap(jsonobj)
	require.NoError(err, msgAndArgs...)

	require.Contains(jsonmap, field, msgAndArgs...)
	require.Nil(jsonmap[field], msgAndArgs...)
}

// RequireJSONFieldNotNil field of jsonobj should not be nil
func RequireJSONFieldNotNil(t *testing.T, jsonobj interface{}, field string, msgAndArgs ...interface{}) {
	require := require.New(t)
	require.NotNil(require)

	jsonmap, err := UnmarshalJSONMap(jsonobj)
	require.NoError(err, msgAndArgs...)

	require.Contains(jsonmap, field, msgAndArgs...)
	require.NotNil(jsonmap[field], msgAndArgs...)
}
