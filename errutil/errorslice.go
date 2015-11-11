package errutil

import (
	"bytes"
	"errors"
)

type ErrorSlice struct {
	ErrorMsg  string   `json:"error,omitempty"`
	ErrorMsgs []string `json:"errors,omitempty"`
	errors    []error
	wrapper   *ErrorWrapper
}

func castErrorSlice(err error) *ErrorSlice {
	switch err.(type) {
	case ErrorSlice:
		errorslice := err.(ErrorSlice)
		return &errorslice
	case *ErrorSlice:
		return err.(*ErrorSlice)
	default:
		return nil
	}
}

func NewErrorSlice(errs ...error) (errorslice *ErrorSlice) {
	e := ErrorSlice{}
	for _, err := range errs {
		if err == nil {
			continue
		}
		es := castErrorSlice(err)
		if es == nil {
			e.errors = append(e.errors, err)
		} else {
			e.errors = append(e.errors, es.errors...)
		}
	}
	if len(e.errors) > 0 {
		e.ErrorMsg = e.errors[0].Error()
		for _, err := range e.errors {
			e.ErrorMsgs = append(e.ErrorMsgs, err.Error())
		}
	} else {
		return nil
	}
	return &e
}

func New(text string, errs ...error) error {
	if text != "" {
		errs = append([]error{errors.New(text)}, errs...)
	}
	return NewErrorSlice(errs...)
}

func Error(errs ...error) error {
	return NewErrorSlice(errs...)
}

func (t ErrorSlice) Error() string {
	if len(t.ErrorMsgs) < 1 {
		return ""
	}
	buffer := bytes.NewBufferString(t.ErrorMsgs[0])
	for _, e := range t.ErrorMsgs[1:] {
		buffer.WriteString("\n")
		buffer.WriteString(e)
	}
	return buffer.String()
}

func (t ErrorSlice) ContainFunc(matchFunc func(error) bool) bool {
	for _, err := range t.errors {
		if matchFunc(err) {
			return true
		}
	}
	return false
}
