package errutil

import (
	"bytes"
	"errors"
)

func New(text string, errs ...error) error {
	e := &ErrorSlice{
		Errs: []error{errors.New(text)},
	}
	for _, err := range errs {
		switch err.(type) {
		case *ErrorSlice:
			e.Errs = append(e.Errs, err.(*ErrorSlice).Errs...)
			break
		default:
			e.Errs = append(e.Errs, err)
			break
		}
	}
	return e
}

type ErrorSlice struct {
	Errs []error
}

func (t *ErrorSlice) Error() string {
	if len(t.Errs) < 1 {
		return ""
	}
	buffer := bytes.NewBufferString(t.Errs[0].Error())
	for _, e := range t.Errs[1:] {
		buffer.WriteString("\n")
		buffer.WriteString(e.Error())
	}
	return buffer.String()
}
