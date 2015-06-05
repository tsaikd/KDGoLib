package injectutil

import (
	"reflect"

	"github.com/codegangsta/inject"
)

// Invoke and handle return value for error type
func Invoke(inj inject.Injector, f interface{}) (refvs []reflect.Value, err error) {
	if refvs, err = inj.Invoke(f); err != nil {
		return
	}
	for _, refv := range refvs {
		if refv.IsValid() {
			refvi := refv.Interface()
			switch refvi.(type) {
			case error:
				err = refvi.(error)
				return
			}
		}
	}
	return
}
