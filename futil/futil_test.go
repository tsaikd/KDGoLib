package futil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsExist(t *testing.T) {
	var (
		assert = assert.New(t)
		err    error
		f      *os.File
	)

	assert.False(IsExist(""), "Empty should not exist")

	f, err = ioutil.TempFile("", "futil_test_")
	assert.NoError(err, "Create TempFile")
	err = f.Close()
	assert.NoError(err, "Close TempFile")

	assert.True(IsExist(f.Name()), "TempFile should exist")

	err = os.Remove(f.Name())
	assert.NoError(err, "Remove TempFile")

	assert.False(IsExist(f.Name()), "TempFile should not exist")
}
