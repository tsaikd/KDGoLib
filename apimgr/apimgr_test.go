package apimgr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tsaikd/KDGoLib/apimgr/teststruct"
	"github.com/tsaikd/KDGoLib/apimgr/teststruct/testsubstruct"
)

func Test_nameGenerator(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	manager := NewManager(Manager{})

	name, fullname := nameGenerator(manager, Definition{
		Version: 1,
		Request: teststruct.TestStruct{},
	})
	assert.Equal("Teststruct", name)
	assert.Equal("TeststructV1", fullname)

	name, fullname = nameGenerator(manager, Definition{
		Version: 2,
		Request: testsubstruct.TestSubStruct{},
	})
	assert.Equal("TestsubstructTeststruct", name)
	assert.Equal("TestsubstructTeststructV2", fullname)
}

func Test_patternGenerator(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	manager := NewManager(Manager{})

	pattern := patternGenerator(manager, Definition{
		Version: 1,
		Request: teststruct.TestStruct{},
	})
	assert.Equal("/1/teststruct", pattern)

	pattern = patternGenerator(manager, Definition{
		Version: 2,
		Request: testsubstruct.TestSubStruct{},
	})
	assert.Equal("/2/testsubstruct/teststruct", pattern)
}
