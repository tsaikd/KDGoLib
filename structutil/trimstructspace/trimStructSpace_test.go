package trimstructspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TrimStructSpace(t *testing.T) {
	assert := assert.New(t)

	func() {
		obj := struct {
			Text      string
			Texts     []string
			SubMap    map[string]string
			SubStruct struct {
				SubText  string
				SubTexts []string
			}
		}{
			Text:  " abc def ",
			Texts: []string{" a ", " b "},
			SubMap: map[string]string{
				"key": " value ",
			},
			SubStruct: struct {
				SubText  string
				SubTexts []string
			}{
				SubText:  " ghi ",
				SubTexts: []string{" jkl "},
			},
		}
		err := TrimStructSpace(&obj)
		assert.NoError(err)
		assert.Equal("abc def", obj.Text)
		assert.Equal([]string{"a", "b"}, obj.Texts)
		assert.Equal("value", obj.SubMap["key"])
		assert.Equal("ghi", obj.SubStruct.SubText)
	}()
}
