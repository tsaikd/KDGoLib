package testcase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SortedMap(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	sortedMap := SortedMap{}
	sortedMap.Set("b", 2)
	sortedMap.Set("a", 1)
	sortedMap.Set("c", 3)
	require.True(sortedMap.IsExists("a"))
	require.False(sortedMap.IsExists("d"))

	sortedMap.Sort()
	require.Equal(3, sortedMap.Len())
	require.Equal(1, sortedMap.Get("a"))
	require.Equal(2, sortedMap.Get("b"))
	require.Equal(3, sortedMap.Get("c"))
	require.Equal("a", sortedMap.Keys()[0])
	require.Equal("b", sortedMap.Keys()[1])
	require.Equal("c", sortedMap.Keys()[2])

	count := 0
	sortedMap.Walk(func(name string, element interface{}) bool {
		count++
		return false
	})
	require.Equal(3, count)

	sortedMap.Remove("b")
	require.Equal(2, sortedMap.Len())
	require.Equal(1, sortedMap.Get("a"))
	require.Equal(3, sortedMap.Get("c"))
	require.Equal("a", sortedMap.Keys()[0])
	require.Equal("c", sortedMap.Keys()[1])

	element := sortedMap.Shift()
	require.Equal(1, element)
	element = sortedMap.Shift()
	require.Equal(3, element)
	element = sortedMap.Shift()
	require.Nil(element)
}
