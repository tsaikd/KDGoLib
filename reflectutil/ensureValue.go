package reflectutil

import "reflect"

// ensure value to avoid pointer
func EnsureValue(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Interface, reflect.Ptr:
		return value.Elem()
	default:
		return value
	}
}
