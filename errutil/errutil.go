package errutil

import (
	"bytes"
	"errors"
)

func New(text string, errs ...error) error {
	e := &typeError{
		errs: []error{errors.New(text)},
	}
	for _, err := range errs {
		switch err.(type) {
		case *typeError:
			e.errs = append(e.errs, err.(*typeError).errs...)
			break
		default:
			e.errs = append(e.errs, err)
			break
		}
	}
	return e
}

type typeError struct {
	errs []error
}

func (t *typeError) Error() string {
	if len(t.errs) < 1 {
		return ""
	}
	buffer := bytes.NewBufferString(t.errs[0].Error())
	for _, e := range t.errs[1:] {
		buffer.WriteString("\n")
		buffer.WriteString(e.Error())
	}
	return buffer.String()
}
