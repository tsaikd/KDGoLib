package jsonutil

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_marshalValue_inline_type(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	val, empty, err := marshalValue(reflect.ValueOf(false), reflect.StructField{})
	require.NoError(err)
	require.Equal(false, val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(true), reflect.StructField{})
	require.NoError(err)
	require.Equal(true, val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(int(0)), reflect.StructField{})
	require.NoError(err)
	require.EqualValues(0, val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(int(1)), reflect.StructField{})
	require.NoError(err)
	require.EqualValues(1, val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(uint(0)), reflect.StructField{})
	require.NoError(err)
	require.EqualValues(0, val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(uint(1)), reflect.StructField{})
	require.NoError(err)
	require.EqualValues(1, val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(float64(0.1)), reflect.StructField{})
	require.NoError(err)
	require.EqualValues(0.1, val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(float64(1.1)), reflect.StructField{})
	require.NoError(err)
	require.EqualValues(1.1, val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(""), reflect.StructField{})
	require.NoError(err)
	require.Equal("", val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf("text"), reflect.StructField{})
	require.NoError(err)
	require.Equal("text", val)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf([]interface{}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf([]interface{}{1}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct{}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Bool bool
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "Bool")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Bool bool `json:",omitempty"`
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Bool bool `json:",omitempty"`
	}{
		Bool: true,
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Bool bool `json:",omitdefault" default:"true"`
	}{
		Bool: true,
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Int int64
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "Int")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Int int64 `json:",omitempty"`
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Int int64 `json:",omitempty"`
	}{
		Int: 1,
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "Int")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Int int64 `json:",omitdefault" default:"1"`
	}{
		Int: 1,
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Uint uint64
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "Uint")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Uint uint64 `json:",omitempty"`
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Uint int64 `json:",omitempty"`
	}{
		Uint: 1,
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "Uint")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Uint int64 `json:",omitdefault" default:"1"`
	}{
		Uint: 1,
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Float float64
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "Float")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Float float64 `json:",omitempty"`
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Float float64 `json:",omitempty"`
	}{
		Float: 1.1,
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "Float")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		Float float64 `json:",omitdefault" default:"1.1"`
	}{
		Float: 1.1,
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		String string
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "String")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		String string `json:",omitempty"`
	}{}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		String string `json:",omitempty"`
	}{
		String: "text",
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 1)
	require.Contains(val, "String")
	require.False(empty)

	val, empty, err = marshalValue(reflect.ValueOf(struct {
		String string `json:",omitdefault" default:"text"`
	}{
		String: "text",
	}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)
}

type testIsEmpty struct {
	empty bool
}

var _ IsEmpty = testIsEmpty{}

func (t testIsEmpty) IsEmpty() bool {
	return t.empty
}

func Test_marshalValue_isEmpty(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	val, empty, err := marshalValue(reflect.ValueOf(testIsEmpty{true}), reflect.StructField{})
	require.NoError(err)
	require.Equal(nil, val)
	require.True(empty)

	val, empty, err = marshalValue(reflect.ValueOf(testIsEmpty{false}), reflect.StructField{})
	require.NoError(err)
	require.Len(val, 0)
	require.False(empty)
}

type testJSONMarshaler struct {
	testIsEmpty
}

var _ json.Marshaler = testJSONMarshaler{}

func (t testJSONMarshaler) MarshalJSON() ([]byte, error) {
	if t.empty {
		return jsonNull, nil
	}
	return []byte("test"), nil
}

func Test_marshalValue_jsonMarshaler(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	val, empty, err := marshalValue(reflect.ValueOf(
		testJSONMarshaler{
			testIsEmpty{true},
		},
	), reflect.StructField{})
	require.NoError(err)
	require.Equal(nil, val)
	require.True(empty)

	val, empty, err = marshalValue(reflect.ValueOf(
		testJSONMarshaler{
			testIsEmpty{false},
		},
	), reflect.StructField{})
	require.NoError(err)
	require.Equal([]byte("test"), val)
	require.False(empty)
}

type testAnonymousChild struct {
	ChildString  string
	childPrivate string
}

type testAnonymousParent struct {
	testAnonymousChild
	ParentString  string
	parentPrivate string
}

func Test_marshalValue_anonymous(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	val, empty, err := marshalValue(reflect.ValueOf(
		testAnonymousParent{
			testAnonymousChild: testAnonymousChild{},
		},
	), reflect.StructField{})
	require.NoError(err)
	require.Contains(val, "ParentString")
	require.NotContains(val, "parentPrivate")
	require.Contains(val, "ChildString")
	require.NotContains(val, "childPrivate")
	require.False(empty)
}

type testMarshalHelper struct {
	testAnonymousParent
}

var _ json.Marshaler = testMarshalHelper{}

func (t testMarshalHelper) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

func Test_MarshalHelper(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	_, err := json.Marshal(testMarshalHelper{})
	require.NoError(err)
}
