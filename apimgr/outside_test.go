package apimgr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tsaikd/KDGoLib/apimgr/testapi"

	_ "github.com/tsaikd/KDGoLib/apimgr/testapi/testres1"
	_ "github.com/tsaikd/KDGoLib/apimgr/testapi/testres2"
)

func Test_Manager(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	totalModule := 3

	manager := testapi.Manager
	assert.Equal("github.com/tsaikd/KDGoLib/apimgr/testapi", manager.BasePackage)

	assert.Len(manager.GetMethodPatternMap(), totalModule)

	sorted := manager.GetSortedAPIsByPkgPath()
	assert.Len(sorted, totalModule)
	if len(sorted) == totalModule {
		assert.Equal("/1/testmethod11/testres1", sorted[0].Pattern)
		assert.Equal("/1/testmethod12/testres1", sorted[1].Pattern)
		assert.Equal("/1/testmethod21/testres2", sorted[2].Pattern)
	}
}
