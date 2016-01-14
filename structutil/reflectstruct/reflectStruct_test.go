package reflectstruct

import (
	"bytes"
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

func (t *enum) UnmarshalJSON(b []byte) (err error) {
	return enumFactory.UnmarshalJSON(t, b)
}

func Test_reflectField(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	func() {
		var field string
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(123),
		)
		if !assert.NoError(err) {
			return
		}
		assert.EqualValues("123", field)
	}()

	func() {
		var field int
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf("123"),
		)
		if !assert.NoError(err) {
			return
		}
		assert.EqualValues(123, field)
	}()

	func() {
		var field int64
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(int8(123)),
		)
		if !assert.NoError(err) {
			return
		}
		assert.EqualValues(123, field)
	}()

	func() {
		var field int64
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(float64(123)),
		)
		if !assert.NoError(err) {
			return
		}
		assert.EqualValues(123, field)
	}()

	func() {
		var field float64
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf("123"),
		)
		if !assert.NoError(err) {
			return
		}
		assert.EqualValues(123, field)
	}()

	func() {
		var field float64
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(int8(123)),
		)
		if !assert.NoError(err) {
			return
		}
		assert.EqualValues(123, field)
	}()

	func() {
		var field float64
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(float32(123)),
		)
		if !assert.NoError(err) {
			return
		}
		assert.EqualValues(123, field)
	}()

	func() {
		var field bool
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf("true"),
		)
		if !assert.NoError(err) {
			return
		}
		assert.True(field)
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
		if !assert.NoError(err) {
			return
		}
		assert.Equal("text", field.Str)
		assert.EqualValues(123, field.Int)
		assert.Equal("intext", field.Inherit.InStr)
		if !assert.Len(field.ChildSlice, 1) {
			return
		}
		assert.Equal("text", field.ChildSlice[0].Str)
		assert.EqualValues(123, field.ChildSlice[0].Int)
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
		assert.NoError(err)
		err = reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf(value),
		)
		if !assert.NoError(err) {
			return
		}
		assert.Equal("text", field.Str)
		assert.EqualValues(123, field.Int)
		assert.Equal("intext", field.Inherit.InStr)
		if !assert.Len(field.ChildSlice, 1) {
			return
		}
		assert.Equal("text", field.ChildSlice[0].Str)
		assert.EqualValues(123, field.ChildSlice[0].Int)
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
		if !assert.NoError(err) {
			return
		}
		assert.Equal("text", field.Str)
		assert.EqualValues(123, field.Int)
		assert.Equal("intext", field.Inherit.InStr)
		if !assert.Len(field.ChildSlice, 1) {
			return
		}
		assert.Equal("text", field.ChildSlice[0].Str)
		assert.EqualValues(123, field.ChildSlice[0].Int)
	}()

	func() {
		var field interface{}
		err := reflectField(
			reflect.ValueOf(&field),
			reflect.ValueOf("text"),
		)
		if !assert.NoError(err) {
			return
		}
		assert.Equal("text", field)
	}()
}

func Test_ReflectStruct_empty_nothing(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

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
	assert.NoError(err)
	assert.Equal("abc", obj.Text)
	assert.Empty(obj.Nothing)
	assert.Empty(obj.Empty)
}

func Test_ReflectStruct_slice(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	obj := struct {
		Text    string   `json:"text"`
		Number  int64    `json:"number"`
		Empty   string   `json:"-"`
		Slice   []string `json:"slice"`
		Numbers []int64  `json:"numbers"`
	}{}
	requrl, err := url.Parse("http://localhost/test?text=abcde&number=123&slice=a&slice=b&numbers=1&numbers=2")
	assert.NoError(err)

	err = ReflectStruct(&obj, requrl.Query())
	assert.NoError(err)
	assert.Equal("abcde", obj.Text)
	assert.Equal(int64(123), obj.Number)
	assert.Empty(obj.Empty)
	assert.Equal([]string{"a", "b"}, obj.Slice)
	assert.Equal([]int64{1, 2}, obj.Numbers)
}

func Test_ReflectStruct_pointer(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	obj := struct {
		PtrStr  string `json:"ptrstr"`
		Str2Int int    `json:"str2int"`
	}{}
	ptrstr := "ptr"
	err := ReflectStruct(&obj, map[string]interface{}{
		"ptrstr":  &ptrstr,
		"str2int": "123",
	})
	assert.NoError(err)
	assert.Equal("ptr", obj.PtrStr)
	assert.Equal(123, obj.Str2Int)
}

func Test_ReflectStruct_type(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	obj := struct {
		Str            fakeString `json:"str"`
		IntFromString  int64      `json:"intFromString"`
		Bool           bool       `json:"bool"`
		BoolFromString bool       `json:"boolFromString"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"str":            "fakestring",
		"intFromString":  "123",
		"bool":           true,
		"boolFromString": "true",
	})
	assert.NoError(err)
	assert.Equal(fakeString("fakestring"), obj.Str)
	assert.EqualValues(123, obj.IntFromString)
	assert.True(obj.Bool)
	assert.True(obj.BoolFromString)
}

func Test_ReflectStruct_struct2struct(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

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
	assert.NoError(err)
}

type childStruct struct {
	ChildString string `json:"childstring"`
}

func Test_ReflectStruct_tag_inherit(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	obj := struct {
		childStruct `reflect:"inherit"`
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
	assert.NoError(err)
	assert.Equal("abc", obj.Text)
	assert.Equal("def", obj.ChildString)
	assert.Equal("ghi", obj.StructField.FieldText)
}

func Test_ReflectStruct_tag_inherit_pointer(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	obj := struct {
		*childStruct `reflect:"inherit"`
		Text         string `json:"text"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"text":        "abc",
		"childstring": "def",
	})
	assert.NoError(err)
	assert.Equal("abc", obj.Text)
	assert.Equal("def", obj.ChildString)
}

func Test_ReflectStruct_enum(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	obj := struct {
		EnumFromString enum `json:"enumFromString"`
		EnumFromBytes  enum `json:"enumFromBytes"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"enumFromString": "test",
		"enumFromBytes":  []byte(`"test"`),
	})
	assert.NoError(err)
	assert.Equal(enumTest, obj.EnumFromString)
	assert.Equal(enumTest, obj.EnumFromBytes)
}
