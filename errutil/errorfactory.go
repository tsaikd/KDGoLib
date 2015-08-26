package errutil

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type ErrorWrapper struct {
	errtext string
	debug   bool
}

func (t *ErrorWrapper) ErrorText(errtext string) *ErrorWrapper {
	t.errtext = errtext
	return t
}

func (t *ErrorWrapper) Debug(debug bool) *ErrorWrapper {
	t.debug = debug
	return t
}

func (t *ErrorWrapper) New(err error, params ...interface{}) error {
	var errtext string
	if t.debug {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}

		file = filepath.Base(file)
		filelinetext := fmt.Sprintf("%s:%d", file, line)

		if t.errtext == "" {
			errtext = filelinetext
		} else {
			errtext = filelinetext + " " + t.errtext
		}
	} else {
		errtext = t.errtext
	}
	curerr := fmt.Errorf(errtext, params...)
	errorslice := NewErrorSlice(curerr, err)
	errorslice.wrapper = t
	return errorslice
}

func (t *ErrorWrapper) Match(err error) bool {
	errorslice := castErrorSlice(err)
	if errorslice == nil {
		return false
	}

	return errorslice.wrapper == t
}

func ErrorFactory(errtext string) *ErrorWrapper {
	return &ErrorWrapper{
		errtext: errtext,
		debug:   false,
	}
}

func ErrorFactoryDebug(errtext string) *ErrorWrapper {
	return &ErrorWrapper{
		errtext: errtext,
		debug:   true,
	}
}
