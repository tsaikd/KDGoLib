package errutil

import "strconv"

// ErrorJSON is a helper struct for display error
type ErrorJSON struct {
	ErrorPath string   `json:"errorpath,omitempty"`
	ErrorMsg  string   `json:"error,omitempty"`
	ErrorMsgs []string `json:"errors,omitempty"`
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
	if err := WalkErrors(errobj, func(errcomp ErrorObject) (stop bool, walkerr error) {
		errors = append(errors, errcomp.Error())
		return false, nil
	}); err != nil {
		return nil
	}

	return &ErrorJSON{
		ErrorPath: errobj.PackageName() + "/" + errobj.FileName() + ":" + strconv.Itoa(errobj.Line()),
		ErrorMsg:  errobj.Error(),
		ErrorMsgs: errors,
	}
}

func (t *ErrorJSON) Error() string {
	if t == nil {
		return ""
	}
	return t.ErrorMsg
}

var _ error = (*ErrorJSON)(nil)
