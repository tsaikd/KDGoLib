package testcase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Queue(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	queue := QueueType{}
	require.Nil(queue.Last())
	require.Nil(queue.Shift())

	queue.Set("b", 2)
	queue.Set("a", 1)
	queue.Set("c", 3)
	require.Equal(3, queue.Last())
	require.True(queue.IsExists("a"))
	require.False(queue.IsExists("d"))

	count := 0
	queue.Walk(func(name string, element interface{}) bool {
		count++
		return false
	})
	require.Equal(3, count)

	queue.Remove("b")
	require.Equal(2, queue.Len())
	require.Equal(1, queue.Get("a"))
	require.Equal(3, queue.Get("c"))
	require.Equal("a", queue.Keys()[0])
	require.Equal("c", queue.Keys()[1])

	element := queue.Shift()
	require.Equal(1, element)
	element = queue.Shift()
	require.Equal(3, element)
	element = queue.Shift()
	require.Nil(element)
}
