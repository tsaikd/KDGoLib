package errutil

import "github.com/tsaikd/KDGoLib/runtimecaller"

// DefaultRuntimeCallerFilter use for filter error stack info
var DefaultRuntimeCallerFilter = []runtimecaller.Filter{}

func init() {
	DefaultRuntimeCallerFilter = append(runtimecaller.FilterCommons, RuntimeCallerFilterStopErrutilPackage)
}

// AddRuntimeCallerFilter add filters to DefaultRuntimeCallerFilter for RuntimeCaller()
func AddRuntimeCallerFilter(filters ...runtimecaller.Filter) {
	DefaultRuntimeCallerFilter = append(DefaultRuntimeCallerFilter, filters...)
}

// RuntimeCallerFilterStopErrutilPackage filter CallInfo to stop after reach KDGoLib/errutil package
func RuntimeCallerFilterStopErrutilPackage(callinfo runtimecaller.CallInfo) (valid bool, stop bool) {
	if callinfo.PackageName == "github.com/tsaikd/KDGoLib/errutil" {
		return false, true
	}
	return true, false
}

// RuntimeCaller wrap runtimecaller.GetByFilters() with DefaultRuntimeCallerFilter
func RuntimeCaller(skip int, extraFilters ...runtimecaller.Filter) (callinfo runtimecaller.CallInfo, ok bool) {
	filters := append(DefaultRuntimeCallerFilter, extraFilters...)
	return runtimecaller.GetByFilters(skip+1, filters...)
}
