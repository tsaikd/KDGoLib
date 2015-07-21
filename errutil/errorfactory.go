package errutil

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type ErrorWrapper func(err error, params ...interface{}) error

func ErrorFactory(errtext string) ErrorWrapper {
	return func(err error, params ...interface{}) error {
		curerr := fmt.Errorf(errtext, params...)
		return NewErrorSlice(curerr, err)
	}
}

func ErrorFactoryDebug(errtext string) ErrorWrapper {
	return func(err error, params ...interface{}) error {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}

		file = filepath.Base(file)
		filelinetext := fmt.Sprintf("%s:%d", file, line)

		if errtext == "" {
			return New(filelinetext, err)
		} else {
			curerr := fmt.Errorf(filelinetext+" "+errtext, params...)
			return NewErrorSlice(curerr, err)
		}
	}
}
