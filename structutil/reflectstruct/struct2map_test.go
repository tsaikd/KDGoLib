package reflectstruct

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ReflectStruct_struct2map(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := map[string]interface{}{}
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

	require.Equal(obj2.Text, obj["text"])
	require.Equal(obj2.Number, obj["number"])
	require.NotContains(obj, "Empty")
	require.NotContains(obj, "-")
	require.Equal(obj2.Slice, obj["slice"])
	require.Equal(obj2.Numbers, obj["numbers"])
}

func Test_ReflectStruct_struct2map_inherit(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := map[string]interface{}{}
	obj2 := struct {
		ChildStruct `reflect:"inherit"`
		Text        string `json:"text"`
	}{
		ChildStruct: ChildStruct{
			ChildString: "child",
		},
		Text: "abc",
	}
	err := ReflectStruct(&obj, obj2)
	require.NoError(err)

	require.Equal(obj2.Text, obj["text"])
	require.Equal("child", obj["childstring"])
}
