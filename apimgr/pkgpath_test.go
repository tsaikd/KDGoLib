package apimgr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getPackagePath(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	pkgpath := getPackagePath(reflect.ValueOf("github.com/tsaikd/KDGoLib/apimgr"))
	assert.Equal("github.com/tsaikd/KDGoLib/apimgr", pkgpath)

	pkgpath = getPackagePath(reflect.ValueOf(Manager{}))
	assert.Equal("github.com/tsaikd/KDGoLib/apimgr", pkgpath)

	pkgpath = getPackagePath(reflect.ValueOf(&Manager{}))
	assert.Equal("github.com/tsaikd/KDGoLib/apimgr", pkgpath)
}
