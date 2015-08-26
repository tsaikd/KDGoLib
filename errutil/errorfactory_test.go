package errutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrorFactory(t *testing.T) {
	assert := assert.New(t)

	errfac1 := ErrorFactory("factory error 1")
	err1 := errfac1.New(nil)
	assert.Error(err1)
	assert.Equal("factory error 1", err1.Error())
	assert.True(errfac1.Match(err1))

	errfac2 := ErrorFactory("factory error 2 with param int %d")
	err2 := errfac2.New(nil, 123)
	assert.Error(err2)
	assert.Equal("factory error 2 with param int 123", err2.Error())
	assert.True(errfac2.Match(err2))
	assert.False(errfac2.Match(err1))

	errfacdebug := ErrorFactoryDebug("debug source")
	errdebug := errfacdebug.New(nil)
	assert.Error(errdebug)
	assert.Equal("errorfactory_test.go:26 debug source", errdebug.Error())
}
