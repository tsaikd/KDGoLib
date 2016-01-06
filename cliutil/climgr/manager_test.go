package climgr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	manager := NewManager()

	err := manager.Add(Definition{
		Name: "A",
	})
	assert.NoError(err)

	modules := manager.GetModules()
	assert.Len(modules, 1)
	assert.Equal("A", modules[0].Name)

	err = manager.Add(Definition{
		Name: "B",
	})
	assert.NoError(err)

	modules = manager.GetModules()
	assert.Len(modules, 2)
	assert.Equal("A", modules[0].Name)
	assert.Equal("B", modules[1].Name)

	err = manager.Add(Definition{
		Name:     "C",
		Priority: -1,
	})
	assert.NoError(err)

	modules = manager.GetModules()
	assert.Len(modules, 3)
	assert.Equal("C", modules[0].Name)
	assert.Equal("A", modules[1].Name)
	assert.Equal("B", modules[2].Name)
}
