package errutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrorFactory(t *testing.T) {
	assert := assert.New(t)

	errfac1 := ErrorFactory("factory error 1")
	err1 := errfac1(nil)
	assert.Error(err1)
	assert.Equal("factory error 1", err1.Error())

	errfac2 := ErrorFactory("factory error 2 with param int %d")
	err2 := errfac2(nil, 123)
	assert.Error(err2)
	assert.Equal("factory error 2 with param int 123", err2.Error())

	errfacdebug := ErrorFactoryDebug("debug source")
	errdebug := errfacdebug(nil)
	assert.Error(errdebug)
	assert.Equal("errorfactory_test.go:23 debug source", errdebug.Error())
}
