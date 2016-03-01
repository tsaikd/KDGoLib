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
	if err == nil {
		return nil
	}

	errobj := castErrorObject(nil, 1, err)
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
		ErrorPath: errobj.Path() + ":" + strconv.Itoa(errobj.Line()),
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
