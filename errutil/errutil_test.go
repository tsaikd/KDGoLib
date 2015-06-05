package errutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	var (
		assert = assert.New(t)
	)

	err1 := New("err1")
	assert.Error(err1)
	assert.Equal("err1", err1.Error())

	err2 := New("err2", err1)
	assert.Error(err2)
	assert.Equal("err2\nerr1", err2.Error())

	err3 := New("err3", err2)
	assert.Error(err3)
	assert.Equal("err3\nerr2\nerr1", err3.Error())

	err4 := New("err4", err2, err1)
	assert.Error(err4)
	assert.Equal("err4\nerr2\nerr1\nerr1", err4.Error())
}
