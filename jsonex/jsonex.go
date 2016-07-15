package jsonex

import (
	"reflect"
	"strconv"
)

const ignoreTag = "-"

// IsEmpty return true if v is empty for supported types
func IsEmpty(v interface{}) bool {
	switch val := v.(type) {
	case reflect.Value:
		return isEmptyValue(val)
	case *reflect.Value:
		if val == nil {
			return true
		}
		return isEmptyValue(*val)
	default:
		return isEmptyValue(reflect.ValueOf(v))
	}
}

func isDefaultValue(v reflect.Value, defaultTag string) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		def := false
		if defaultTag != "" {
			var err error
			if def, err = strconv.ParseBool(defaultTag); err != nil {
				panic(err)
			}
		}
		return v.Bool() == def
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		def := int64(0)
		if defaultTag != "" {
			var err error
			if def, err = strconv.ParseInt(defaultTag, 0, 64); err != nil {
				panic(err)
			}
		}
		return v.Int() == def
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		def := uint64(0)
		if defaultTag != "" {
			var err error
			if def, err = strconv.ParseUint(defaultTag, 0, 64); err != nil {
				panic(err)
			}
		}
		return v.Uint() == def
	case reflect.Float32, reflect.Float64:
		def := float64(0)
		if defaultTag != "" {
			var err error
			if def, err = strconv.ParseFloat(defaultTag, 64); err != nil {
				panic(err)
			}
		}
		return v.Float() == def
	case reflect.String:
		return v.String() == defaultTag
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return isDefaultValue(v.Elem(), defaultTag)
	case reflect.Struct:
		t := v.Type()
		for i, n := 0, t.NumField(); i < n; i++ {
			elemVal := v.Field(i)
			elemField := t.Field(i)
			if !elemField.Anonymous && !elemVal.CanInterface() {
				continue
			}

			tag := elemField.Tag.Get("json")
			if tag == ignoreTag {
				continue
			}

			if isDefaultValue(elemVal, elemField.Tag.Get("default")) {
				continue
			}

			return false
		}
		return true
	}
	return false
}
