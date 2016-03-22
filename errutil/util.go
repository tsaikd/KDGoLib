package errutil

import (
	"runtime"
	"strings"
)

// RuntimeCallerFilter use to filter runtime.Caller result
type RuntimeCallerFilter func(file string, line int) bool

// RuntimeCaller wrap runtime.Caller(), find until go source file
func RuntimeCaller(skip int, filters ...RuntimeCallerFilter) (file string, line int, ok bool) {
	filters = append([]RuntimeCallerFilter{
		func(file string, line int) bool {
			return strings.HasSuffix(file, ".go")
		},
	}, filters...)
	return RuntimeCallerCustomFilter(skip+1, filters...)
}

func filterAll(file string, line int, filters ...RuntimeCallerFilter) bool {
	for _, filter := range filters {
		if !filter(file, line) {
			return false
		}
	}
	return true
}

// RuntimeCallerCustomFilter wrap runtime.Caller(), find until all filters match
func RuntimeCallerCustomFilter(skip int, filters ...RuntimeCallerFilter) (file string, line int, ok bool) {
	skip++

	for {
		_, file, line, ok = runtime.Caller(skip)
		if !ok {
			return
		}

		if filterAll(file, line, filters...) {
			return
		}

		skip++
	}
}

// ContainErrorFunc check error contain error by custom equalFunc
func ContainErrorFunc(err error, equalFunc func(error) bool) bool {
	errobj := castErrorObject(nil, 1, err)
	contain := false

	if walkerr := WalkErrors(errobj, func(errcomp ErrorObject) (stop bool, walkerr error) {
		if equalFunc(errcomp) {
			contain = true
			return true, nil
		}
		return false, nil
	}); walkerr != nil {
		return false
	}

	return contain
}
