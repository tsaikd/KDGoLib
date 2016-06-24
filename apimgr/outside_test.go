package apimgr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/apimgr/testapi"

	_ "github.com/tsaikd/KDGoLib/apimgr/testapi/testres1"
	_ "github.com/tsaikd/KDGoLib/apimgr/testapi/testres2"
)

func Test_Manager(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	totalModule := 3

	manager := testapi.Manager
	require.Contains(manager.BasePackage, "github.com/tsaikd/KDGoLib/apimgr/testapi")

	require.Len(manager.GetMethodPatternMap(), totalModule)

	sorted := manager.GetSortedAPIsByPkgPath()
	if assert.Len(sorted, totalModule) {
		require.Equal("/1/testmethod11/testres1", sorted[0].Pattern)
		require.Equal("/1/testmethod12/testres1", sorted[1].Pattern)
		require.Equal("/1/testmethod21/testres2", sorted[2].Pattern)
	}
}
