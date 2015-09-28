package reflectstruct

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeString string

func Test_ReflectStruct(t *testing.T) {
	assert := assert.New(t)

	func() {
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
		assert.Empty(obj.Empty)
	}()

	func() {
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
		assert.Equal(123, obj.Number)
		assert.Empty(obj.Empty)
		assert.Equal([]string{"a", "b"}, obj.Slice)
		assert.Equal([]int64{1, 2}, obj.Numbers)
	}()

	func() {
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
	}()

	func() {
		obj := struct {
			Str fakeString `json:"str"`
		}{}
		err := ReflectStruct(&obj, map[string]interface{}{
			"str": "fakestring",
		})
		assert.NoError(err)
		assert.Equal("fakestring", obj.Str)
	}()

	func() {
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
	}()
}
