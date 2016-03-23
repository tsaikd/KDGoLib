package errutil

import (
	"errors"
	"path/filepath"
)

// errors
var (
	ErrorWalkLoop = NewFactory("detect error component parents loop when walking")
)

// New return a new ErrorObject object
func New(text string, errs ...error) ErrorObject {
	if text != "" {
		errs = append([]error{errors.New(text)}, errs...)
	}
	return NewErrorsSkip(1, errs...)
}

// NewErrors return ErrorObject that contains all input errors
func NewErrors(errs ...error) ErrorObject {
	return NewErrorsSkip(1, errs...)
}

// NewErrorsSkip return ErrorObject, skip function call
func NewErrorsSkip(skip int, errs ...error) ErrorObject {
	var errcomp ErrorObject
	var errtmp ErrorObject
	for i, size := 0, len(errs); i < size; i++ {
		errtmp = castErrorObject(nil, skip+1, errs[i])
		if errtmp == nil {
			continue
		}

		if errcomp == nil {
			errcomp = errtmp
			continue
		}

		if err := AddParent(errcomp, errtmp); err != nil {
			panic(err)
		}
	}
	return errcomp
}

// ErrorObject is a rich error interface
type ErrorObject interface {
	Path() string
	Filename() string
	Line() int
	Error() string
	Factory() ErrorFactory
	Parent() ErrorObject
	SetParent(parent ErrorObject) ErrorObject
}

type errorObject struct {
	path     string
	filename string
	line     int
	errtext  string
	factory  ErrorFactory
	parent   ErrorObject
}

func castErrorObject(factory ErrorFactory, skip int, err error) ErrorObject {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case errorObject:
		res := err.(errorObject)
		return &res
	case *errorObject:
		return err.(*errorObject)
	case ErrorObject:
		return err.(ErrorObject)
	default:
		file, line, _ := RuntimeCaller(skip + 1)
		filename := filepath.Base(file)
		return &errorObject{
			path:     file,
			filename: filename,
			line:     line,
			errtext:  err.Error(),
			factory:  factory,
		}
	}
}

func (t *errorObject) Path() string {
	if t == nil {
		return ""
	}
	return t.path
}

func (t *errorObject) Filename() string {
	if t == nil {
		return ""
	}
	return t.filename
}

func (t *errorObject) Line() int {
	if t == nil {
		return 0
	}
	return t.line
}

func (t errorObject) Error() string {
	return t.errtext
}

func (t *errorObject) Factory() ErrorFactory {
	if t == nil {
		return nil
	}
	return t.factory
}

func (t *errorObject) Parent() ErrorObject {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *errorObject) SetParent(parent ErrorObject) ErrorObject {
	if t == nil {
		return nil
	}
	if t == parent {
		return t
	}
	t.parent = parent
	return t
}

func (t *errorObject) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

var _ ErrorObject = (*errorObject)(nil)
