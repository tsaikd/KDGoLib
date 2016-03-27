package errutil

import (
	"strconv"
	"strings"
)

// ErrorJSON is a helper struct for display error
type ErrorJSON struct {
	ErrorPath      string          `json:"errorpath,omitempty"`
	ErrorMsg       string          `json:"error,omitempty"`
	ErrorMsgs      []string        `json:"errors,omitempty"`
	ErrorFactories map[string]bool `json:"errfac,omitempty"`
}

// NewJSON create ErrorJSON
func NewJSON(err error) *ErrorJSON {
	return newJSON(1, err)
}

func newJSON(skip int, err error) *ErrorJSON {
	errobj := castErrorObject(nil, skip+1, err)
	if errobj == nil {
		return nil
	}

	errors := []string{}
	facs := map[string]bool{}
	if err := WalkErrors(errobj, func(errcomp ErrorObject) (stop bool, walkerr error) {
		errors = append(errors, errcomp.Error())
		factory := errcomp.Factory()
		if factory != nil {
			if !strings.Contains(factory.Name(), "->") {
				facs[factory.Name()] = true
			}
		}
		return false, nil
	}); err != nil {
		return nil
	}

	return &ErrorJSON{
		ErrorPath:      errobj.PackageName() + "/" + errobj.FileName() + ":" + strconv.Itoa(errobj.Line()),
		ErrorMsg:       errobj.Error(),
		ErrorMsgs:      errors,
		ErrorFactories: facs,
	}
}

func (t *ErrorJSON) Error() string {
	if t == nil {
		return ""
	}
	return t.ErrorMsg
}

var _ error = (*ErrorJSON)(nil)
