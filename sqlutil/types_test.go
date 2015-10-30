package sqlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	testStringSlice := SQLStringSlice{"abc", "def"}
	sqlvalue, err := SQLValueStringSlice(testStringSlice)
	assert.NoError(err)
	assert.Equal(`{"abc","def"}`, sqlvalue)

	sqlvalue, err = SQLValueStringSlice(&testStringSlice)
	assert.NoError(err)
	assert.Equal(`{"abc","def"}`, sqlvalue)

	sqlvalue, err = SQLValueStringSlice(nil)
	assert.NoError(err)
	assert.Nil(sqlvalue)
}
