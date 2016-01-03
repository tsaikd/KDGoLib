package injectutil

import (
	"reflect"

	"github.com/codegangsta/inject"
)

// CheckErrorValues return error if refvs contains error
func CheckErrorValues(refvs []reflect.Value) (err error) {
	for _, refv := range refvs {
		if refv.IsValid() {
			refvi := refv.Interface()
			switch refvi.(type) {
			case error:
				return refvi.(error)
			}
		}
	}
	return
}

// Invoke and handle return value for error type
func Invoke(inj inject.Injector, f interface{}) (refvs []reflect.Value, err error) {
	if refvs, err = inj.Invoke(f); err != nil {
		return
	}
	err = CheckErrorValues(refvs)
	return
}
