package reflectstruct

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ReflectStruct_map2map(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	obj := map[string]interface{}{}
	obj2 := map[string]interface{}{
		"Text":    "abc",
		"Number":  123,
		"Slice":   []string{"a", "b", "c"},
		"Numbers": []int64{1, 2, 3},
	}
	err := ReflectStruct(&obj, obj2)
	require.NoError(err)

	require.Equal(obj2["Text"], obj["Text"])
	require.Equal(obj2["Number"], obj["Number"])
	require.Equal(obj2["Slice"], obj["Slice"])
	require.Equal(obj2["Numbers"], obj["Numbers"])

	obj = map[string]interface{}{}
	obj2 = nil
	err = ReflectStruct(&obj, obj2)
	require.NoError(err)
	require.NotNil(obj)
}
