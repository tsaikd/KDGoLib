package runtimecaller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeCaller(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	// clean FilterCommons for test in this package
	commons := FilterCommons
	FilterCommons = []Filter{}
	defer func() {
		FilterCommons = commons
	}()

	callinfo, ok := GetByFilters(0)
	require.True(ok)
	require.NotZero(callinfo.PC())
	require.NotZero(callinfo.FilePath())
	require.NotZero(callinfo.Line())
	require.NotNil(callinfo.PCFunc())
	require.NotZero(callinfo.PackageName())
	require.NotZero(callinfo.FileDir())
	require.NotZero(callinfo.FileName())
	require.NotZero(callinfo.FuncName())

	callinfos := ListByFilters(0)
	fullstacklen := len(callinfos)
	require.True(fullstacklen > 1)

	callinfos = ListByFilters(0, FilterOnlyGoSource)
	require.True(fullstacklen >= len(callinfos))

	callinfos = ListByFilters(0, FilterOnlyGoSource, FilterStopRuntimeCallerPackage)
	require.True(fullstacklen >= len(callinfos))
}

func Test_retrieveCallInfo(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	require.NotPanics(func() {
		ok := true
		tried := false
		for i := 1; ok; i++ {
			_, ok = retrieveCallInfo(i)
			if ok {
				tried = true
			}
		}
		require.True(tried)
	})
}
