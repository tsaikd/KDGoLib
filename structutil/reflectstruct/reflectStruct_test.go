package reflectstruct

import (
	"bytes"
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/enumutil"
)

type fakeString string

type enum int8

const (
	enumTest enum = 1 + iota
)

var enumFactory = enumutil.NewEnumFactory().
	Add(enumTest, "test").
	Build()

func (t enum) String() string {
	return enumFactory.String(t)
}

func (t *enum) UnmarshalJSON(b []byte) (err error) {
	return enumFactory.UnmarshalJSON(t, b)
}

func Test_reflectField_simple_value(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	var tmpstr string
	var tmpint int
	var tmpint64 int64
	var tmpfloat64 float64
	var tmpbool bool
	var tmpenum enum
	var tmptime time.Time
	var tmpvar interface{}

	now := time.Now()

	// normal usage
	for _, data := range []struct {
		dest     reflect.Value
		value    reflect.Value
		expected interface{}
	}{
		{reflect.ValueOf(&tmpstr), reflect.ValueOf(false), "false"},
		{reflect.ValueOf(&tmpstr), reflect.ValueOf(true), "true"},
		{reflect.ValueOf(&tmpstr), reflect.ValueOf(int(123)), "123"},
		{reflect.ValueOf(&tmpstr), reflect.ValueOf(float64(1.23)), "1.23"},
		{reflect.ValueOf(&tmpint), reflect.ValueOf(false), int(0)},
		{reflect.ValueOf(&tmpint), reflect.ValueOf(true), int(1)},
		{reflect.ValueOf(&tmpint), reflect.ValueOf("123"), int(123)},
		{reflect.ValueOf(&tmpint64), reflect.ValueOf(int8(123)), int64(123)},
		{reflect.ValueOf(&tmpint64), reflect.ValueOf(float64(123)), int64(123)},
		{reflect.ValueOf(&tmpfloat64), reflect.ValueOf(false), float64(0)},
		{reflect.ValueOf(&tmpfloat64), reflect.ValueOf(true), float64(1)},
		{reflect.ValueOf(&tmpfloat64), reflect.ValueOf("123"), float64(123)},
		{reflect.ValueOf(&tmpfloat64), reflect.ValueOf(int8(123)), float64(123)},
		{reflect.ValueOf(&tmpfloat64), reflect.ValueOf(float32(123)), float64(123)},
		{reflect.ValueOf(&tmpbool), reflect.ValueOf("true"), true},
		{reflect.ValueOf(&tmpbool), reflect.ValueOf(int(123)), true},
		{reflect.ValueOf(&tmpbool), reflect.ValueOf(float32(123)), true},
		{reflect.ValueOf(&tmpenum), reflect.ValueOf("test"), enumTest},
		{reflect.ValueOf(&tmpenum), reflect.ValueOf(2), enum(2)},
		{reflect.ValueOf(&tmptime), reflect.ValueOf(now), now},
		{reflect.ValueOf(&tmpvar), reflect.ValueOf("text"), "text"},
	} {
		require.NoError(reflectField(data.dest, data.value))
		require.Equal(data.expected, data.dest.Elem().Interface())
	}

	// error usage
	for _, data := range []struct {
		dest  reflect.Value
		value reflect.Value
	}{
		{reflect.ValueOf(&tmpbool), reflect.ValueOf("abc")},
		{reflect.ValueOf(&tmpint), reflect.ValueOf("abc")},
		{reflect.ValueOf(&tmpfloat64), reflect.ValueOf("abc")},
	} {
		require.Error(reflectField(data.dest, data.value))
	}
}

func Test_reflectField_struct(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	func() {
		var field struct {
			Str     string `json:"str"`
			Int     int    `json:"int"`
			Inherit struct {
				InStr string `json:"instr"`
			} `reflect:"inherit"`
			ChildSlice []struct {
				Str string `json:"str"`
				Int int    `json:"int"`
			} `json:"childslice"`
		}
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(map[string]interface{}{
				"str":   "text",
				"int":   123,
				"instr": "intext",
				"childslice": []map[string]interface{}{
					map[string]interface{}{
						"str": "text",
						"int": 123,
					},
				},
			}),
		)
		require.NoError(err)
		require.Equal("text", field.Str)
		require.EqualValues(123, field.Int)
		require.Equal("intext", field.Inherit.InStr)
		require.Len(field.ChildSlice, 1)
		require.Equal("text", field.ChildSlice[0].Str)
		require.EqualValues(123, field.ChildSlice[0].Int)
	}()

	func() {
		var field struct {
			Str     string `json:"str"`
			Int     int    `json:"int"`
			Inherit struct {
				InStr string `json:"instr"`
			} `reflect:"inherit"`
			ChildSlice []struct {
				Str string `json:"str"`
				Int int    `json:"int"`
			} `json:"childslice"`
		}
		var value interface{}
		err := json.NewDecoder(bytes.NewBufferString(`
			{
				"str": "text",
				"int": 123,
				"instr": "intext",
				"childslice": [
					{
						"str": "text",
						"int": 123
					}
				]
			}
		`)).Decode(&value)
		require.NoError(err)
		err = reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(value),
		)
		require.NoError(err)
		require.Equal("text", field.Str)
		require.EqualValues(123, field.Int)
		require.Equal("intext", field.Inherit.InStr)
		require.Len(field.ChildSlice, 1)
		require.Equal("text", field.ChildSlice[0].Str)
		require.EqualValues(123, field.ChildSlice[0].Int)
	}()

	func() {
		var field struct {
			Str     string `json:"str"`
			Int     int    `json:"int"`
			Inherit struct {
				InStr string `json:"instr"`
			} `reflect:"inherit"`
			ChildSlice []struct {
				Str string `json:"str"`
				Int int    `json:"int"`
			} `json:"childslice"`
		}
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(struct {
				ValStr        string `json:"str"`
				ValInt        int    `json:"int"`
				ValInStr      string `json:"instr"`
				ValChildSlice []struct {
					ValStr string `json:"str"`
					ValInt int    `json:"int"`
				} `json:"childslice"`
			}{
				ValStr:   "text",
				ValInt:   123,
				ValInStr: "intext",
				ValChildSlice: []struct {
					ValStr string `json:"str"`
					ValInt int    `json:"int"`
				}{
					struct {
						ValStr string `json:"str"`
						ValInt int    `json:"int"`
					}{
						ValStr: "text",
						ValInt: 123,
					},
				},
			}),
		)
		require.NoError(err)
		require.Equal("text", field.Str)
		require.EqualValues(123, field.Int)
		require.Equal("intext", field.Inherit.InStr)
		require.Len(field.ChildSlice, 1)
		require.Equal("text", field.ChildSlice[0].Str)
		require.EqualValues(123, field.ChildSlice[0].Int)
	}()
}

func Test_ReflectStruct_nil(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err := ReflectStruct(nil, nil)
	require.NoError(err)
}

func Test_ReflectStruct_empty_nothing(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := struct {
		Text    string `json:"text"`
		Nothing string
		Empty   string `json:"-"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"text":    "abc",
		"nothing": 123,
		"empty":   "null",
	})
	require.NoError(err)
	require.Equal("abc", obj.Text)
	require.Empty(obj.Nothing)
	require.Empty(obj.Empty)
}

func Test_ReflectStruct_slice(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := struct {
		Text    string   `json:"text"`
		Number  int64    `json:"number"`
		Empty   string   `json:"-"`
		Slice   []string `json:"slice"`
		Numbers []int64  `json:"numbers"`
	}{}
	requrl, err := url.Parse("http://localhost/test?text=abcde&number=123&slice=a&slice=b&numbers=1&numbers=2")
	require.NoError(err)

	err = ReflectStruct(&obj, requrl.Query())
	require.NoError(err)
	require.Equal("abcde", obj.Text)
	require.Equal(int64(123), obj.Number)
	require.Empty(obj.Empty)
	require.Equal([]string{"a", "b"}, obj.Slice)
	require.Equal([]int64{1, 2}, obj.Numbers)

	obj2 := struct {
		Empty string `json:"empty"`
	}{}
	require.NoError(ReflectStruct(&obj2, map[string]interface{}{
		"empty": []string{},
	}))
	require.Empty(obj2.Empty)
}

func Test_ReflectStruct_pointer(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := struct {
		PtrStr  string `json:"ptrstr"`
		Str2Int int    `json:"str2int"`
	}{}
	ptrstr := "ptr"
	err := ReflectStruct(&obj, map[string]interface{}{
		"ptrstr":  &ptrstr,
		"str2int": "123",
	})
	require.NoError(err)
	require.Equal("ptr", obj.PtrStr)
	require.Equal(123, obj.Str2Int)
}

func Test_ReflectStruct_type(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := struct {
		Str               fakeString `json:"str"`
		IntFromString     int64      `json:"intFromString"`
		Bool              bool       `json:"bool"`
		BoolFromString    bool       `json:"boolFromString"`
		BoolPtr           *bool      `json:"boolPtr"`
		BoolPtrFromString *bool      `json:"boolPtrFromString"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"str":               "fakestring",
		"intFromString":     "123",
		"bool":              true,
		"boolFromString":    "true",
		"boolPtr":           true,
		"boolPtrFromString": "true",
	})
	require.NoError(err)
	require.Equal(fakeString("fakestring"), obj.Str)
	require.EqualValues(123, obj.IntFromString)
	require.True(obj.Bool)
	require.True(obj.BoolFromString)
	require.NotNil(obj.BoolPtr)
	require.True(*obj.BoolPtr)
	require.NotNil(obj.BoolPtrFromString)
	require.True(*obj.BoolPtrFromString)
}

func Test_ReflectStruct_struct2struct(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := struct {
		Text    string   `json:"text"`
		Number  int64    `json:"number"`
		Empty   string   `json:"-"`
		Slice   []string `json:"slice"`
		Numbers []int64  `json:"numbers"`
	}{}
	obj2 := struct {
		Text    string   `json:"text"`
		Number  int64    `json:"number"`
		Empty   string   `json:"-"`
		Slice   []string `json:"slice"`
		Numbers []int64  `json:"numbers"`
	}{
		Text:    "abc",
		Number:  123,
		Slice:   []string{"a", "b", "c"},
		Numbers: []int64{1, 2, 3},
	}
	err := ReflectStruct(&obj, obj2)
	require.NoError(err)
}

type ChildStruct struct {
	ChildString string `json:"childstring"`
}

func Test_ReflectStruct_tag_inherit(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := struct {
		ChildStruct `reflect:"inherit"`
		Text        string `json:"text"`
		StructField struct {
			FieldText string `json:"fieldText"`
		} `json:"structField"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"text":        "abc",
		"childstring": "def",
		"structField": map[string]interface{}{
			"fieldText": "ghi",
		},
	})
	require.NoError(err)
	require.Equal("abc", obj.Text)
	require.Equal("def", obj.ChildString)
	require.Equal("ghi", obj.StructField.FieldText)
}

func Test_ReflectStruct_tag_inherit_pointer(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := struct {
		*ChildStruct `reflect:"inherit"`
		Text         string `json:"text"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"text":        "abc",
		"childstring": "def",
	})
	require.NoError(err)
	require.Equal("abc", obj.Text)
	require.Equal("def", obj.ChildString)
}

func Test_ReflectStruct_enum(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := struct {
		EnumFromString enum `json:"enumFromString"`
		EnumFromBytes  enum `json:"enumFromBytes"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"enumFromString": "test",
		"enumFromBytes":  []byte(`"test"`),
	})
	require.NoError(err)
	require.Equal(enumTest, obj.EnumFromString)
	require.Equal(enumTest, obj.EnumFromBytes)
}
