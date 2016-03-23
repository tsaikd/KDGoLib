package errutil

import (
	"path"
	"runtime"
	"strings"
)

// RuntimeCallerFilter use to filter runtime.Caller result
type RuntimeCallerFilter func(packageName string, fileName string, funcName string, line int) bool

// RuntimeCaller wrap runtime.Caller(), find until go source file
func RuntimeCaller(skip int, filters ...RuntimeCallerFilter) (packageName string, fileName string, funcName string, line int, ok bool) {
	filters = append([]RuntimeCallerFilter{
		func(packageName string, fileName string, funcName string, line int) bool {
			return strings.HasSuffix(fileName, ".go")
		},
	}, filters...)
	return RuntimeCallerCustomFilter(skip+1, filters...)
}

// http://stackoverflow.com/questions/25262754/how-to-get-name-of-current-package-in-go
func retrieveCallInfo(skip int) (_packageName string, _fileName string, _funcName string, _line int, _ok bool) {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return
	}
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return packageName, fileName, funcName, line, true
}

func filterAll(packageName string, fileName string, funcName string, line int, filters ...RuntimeCallerFilter) bool {
	for _, filter := range filters {
		if !filter(packageName, fileName, funcName, line) {
			return false
		}
	}
	return true
}

// RuntimeCallerCustomFilter wrap runtime.Caller(), find until all filters match
func RuntimeCallerCustomFilter(skip int, filters ...RuntimeCallerFilter) (packageName string, fileName string, funcName string, line int, ok bool) {
	skip++

	for {
		packageName, fileName, funcName, line, ok = retrieveCallInfo(skip)
		if !ok {
			return
		}

		if filterAll(packageName, fileName, funcName, line, filters...) {
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
