package runtimecaller

import (
	"path"
	"runtime"
	"strings"
)

// GetByFilters return CallInfo until all filters are valid
func GetByFilters(skip int, filters ...Filter) (callinfo CallInfo, ok bool) {
	skip++

	for {
		if callinfo, ok = retrieveCallInfo(skip); !ok {
			return
		}

		valid, stop := filterAll(callinfo, filters...)
		if valid {
			return callinfo, true
		}
		if stop {
			return callinfo, false
		}

		skip++
	}
}

// ListByFilters return all CallInfo stack for all filters are valid
func ListByFilters(skip int, filters ...Filter) (callinfos []CallInfo) {
	skip++

	for {
		var callinfo CallInfo
		var ok bool

		if callinfo, ok = retrieveCallInfo(skip); !ok {
			return
		}

		valid, stop := filterAll(callinfo, filters...)
		if valid {
			callinfos = append(callinfos, callinfo)
		}
		if stop {
			return
		}

		skip++
	}
}

// http://stackoverflow.com/questions/25262754/how-to-get-name-of-current-package-in-go
func retrieveCallInfo(skip int) (callinfo CallInfo, ok bool) {
	if callinfo.PC, callinfo.FilePath, callinfo.Line, ok = runtime.Caller(skip + 1); !ok {
		return
	}

	callinfo.FileDir, callinfo.FileName = path.Split(callinfo.FilePath)
	callinfo.PCFunc = runtime.FuncForPC(callinfo.PC)

	parts := strings.Split(callinfo.PCFunc.Name(), ".")
	pl := len(parts)
	callinfo.FuncName = parts[pl-1]

	if parts[pl-2][0] == '(' {
		callinfo.FuncName = parts[pl-2] + "." + callinfo.FuncName
		callinfo.PackageName = strings.Join(parts[0:pl-2], ".")
	} else {
		callinfo.PackageName = strings.Join(parts[0:pl-1], ".")
	}

	return callinfo, true
}

func filterAll(callinfo CallInfo, filters ...Filter) (allvalid bool, onestop bool) {
	allvalid = true
	for _, filter := range filters {
		valid, stop := filter(callinfo)
		allvalid = allvalid && valid
		if stop {
			return allvalid, true
		}
		if !allvalid {
			return
		}
	}
	return
}
