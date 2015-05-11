package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	assert := assert.New(t)

	omap := New("k", "d", "o", "r", "d", "e", "r", "e", "d", "m", "a", "p")

	assert.True(omap.Exist("d"))
	assert.False(omap.Exist("z"))

	assert.Len(omap.Slice(), 8)

	data, err := omap.MarshalJSON()
	assert.NoError(err)
	assert.Equal(`["k","d","o","r","e","m","a","p"]`, string(data))
}
