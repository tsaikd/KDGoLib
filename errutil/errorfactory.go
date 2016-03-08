package errutil

import "fmt"

// ErrorFactory is used for create or check ErrorObject
type ErrorFactory interface {
	Error() string

	New(err error, params ...interface{}) ErrorObject
	Match(err error) bool
	In(err error) bool
}

type errorFactory struct {
	errtext string
}

// NewFactory return new NewFactory instance
func NewFactory(errtext string) ErrorFactory {
	return &errorFactory{
		errtext: errtext,
	}
}

func (t errorFactory) Error() string {
	return t.errtext
}

func (t *errorFactory) New(parent error, params ...interface{}) ErrorObject {
	errobj := castErrorObject(t, 1, fmt.Errorf(t.errtext, params...))
	errobj.SetParent(castErrorObject(nil, 1, parent))
	return errobj
}

func (t *errorFactory) Match(err error) bool {
	if t == nil || err == nil {
		return false
	}

	errcomp := castErrorObject(nil, 1, err)
	if errcomp == nil {
		return false
	}

	return errcomp.Factory() == t
}

func (t *errorFactory) In(err error) bool {
	if t == nil || err == nil {
		return false
	}

	exist := false

	if errtmp := WalkErrors(castErrorObject(nil, 1, err), func(errcomp ErrorObject) (stop bool, walkerr error) {
		if errcomp.Factory() == t {
			exist = true
			return true, nil
		}
		return false, nil
	}); errtmp != nil {
		panic(errtmp)
	}

	return exist
}

var _ ErrorFactory = (*errorFactory)(nil)
