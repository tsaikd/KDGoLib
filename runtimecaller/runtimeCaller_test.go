package runtimecaller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RuntimeCaller(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	callinfo, ok := GetByFilters(0)
	require.True(ok)
	require.NotZero(callinfo.PC)
	require.NotZero(callinfo.FilePath)
	require.NotZero(callinfo.Line)
	require.NotNil(callinfo.PCFunc)
	require.NotZero(callinfo.PackageName)
	require.NotZero(callinfo.FileDir)
	require.NotZero(callinfo.FileName)
	require.NotZero(callinfo.FuncName)

	callinfos := ListByFilters(0)
	fullstacklen := len(callinfos)
	require.True(fullstacklen > 1)

	callinfos = ListByFilters(0, FilterOnlyGoSource)
	require.True(fullstacklen >= len(callinfos))

	callinfos = ListByFilters(0, FilterOnlyGoSource, FilterStopRuntimeCallerPackage)
	require.True(fullstacklen >= len(callinfos))

	assertinfo := assert.CallerInfo()
	t.Log(assertinfo)
}
