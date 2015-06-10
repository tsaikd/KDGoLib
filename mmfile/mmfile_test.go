package mmfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadZero(t *testing.T) {
	assert := assert.New(t)

	mf, err := Open("testdata/zero.bin")
	assert.NoError(err)
	assert.Len(mf.Data(), 32)

	data := mf.Data()
	for _, b := range data {
		assert.Equal(byte('\x00'), b)
	}

	err = mf.Close()
	assert.NoError(err)
}

func TestReadHello(t *testing.T) {
	assert := assert.New(t)

	mf, err := Open("testdata/hello.txt")
	assert.NoError(err)

	data := mf.Data()
	assert.Equal([]byte("Hello World"), data[:11])

	err = mf.Close()
	assert.NoError(err)
}
