package errutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

// Trace error stack, output to os.Stderr
func Trace(err error) {
	trace(1, err, os.Stderr)
}

// TraceWriter error stack, output to writer
func TraceWriter(err error, writer io.Writer) {
	trace(1, err, writer)
}

func trace(skip int, err error, writer io.Writer) {
	errjson := newJSON(skip+1, err)
	if errjson == nil {
		return
	}

	data, err := json.Marshal(errjson)
	if err != nil {
		panic(err)
	}
	if len(data) < 1 {
		return
	}

	if writer == nil {
		writer = os.Stderr
	}

	if _, werr := fmt.Fprintln(writer, string(data)); writer != os.Stderr && werr != nil {
		fmt.Fprintln(os.Stderr, werr, string(data))
	}
}
