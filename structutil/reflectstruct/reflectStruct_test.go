package reflectstruct

import (
	"net/url"
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
		Bool           bool       `json:"bool"`
		BoolFromString bool       `json:"boolFromString"`
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"str":            "fakestring",
		"bool":           true,
		"boolFromString": "true",
	})
	assert.NoError(err)
	assert.Equal(fakeString("fakestring"), obj.Str)
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
	}{}
	err := ReflectStruct(&obj, map[string]interface{}{
		"text":        "abc",
		"childstring": "def",
	})
	assert.NoError(err)
	assert.Equal("abc", obj.Text)
	assert.Equal("def", obj.ChildString)
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
