package sqlutil

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tsaikd/KDGoLib/errutil"
)

var (
	ErrorUnsupportedScanType1 = errutil.NewFactory("unsupported scan type: %v")
	ErrorInvalidObjType1      = errutil.NewFactory("invalid obj type: %T")
	ErrorNoValueFound1        = errutil.NewFactory("no value found for key %v")
)

// Set obj to value's JSON representation
func SQLScanJson(obj interface{}, value interface{}) (err error) {
	if value == nil {
		return
	}
	switch value.(type) {
	case []byte:
		return json.Unmarshal(value.([]byte), obj)
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}

// Return obj's JSON representation which implements driver.Value
func SQLValueJson(obj interface{}) (value driver.Value, err error) {
	if obj == nil {
		return
	}
	jsondata, err := json.Marshal(obj)
	if err != nil {
		return
	}
	if bytes.Equal([]byte("{}"), jsondata) {
		return
	}
	value = jsondata
	return
}

func ensureScanValue(obj interface{}) (refval reflect.Value, err error) {
	// Get the value of obj and make sure it's either a pointer or nil
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return rv, ErrorInvalidObjType1.New(nil, obj)
	}
	// So we can get its actual value
	pv := reflect.Indirect(rv)

	return pv, nil
}

// Set obj to value's string representation
func SQLScanString(obj interface{}, value interface{}) (err error) {
	if value == nil {
		return
	}

	pv, err := ensureScanValue(obj)
	if err != nil {
		return
	}

	switch value.(type) {
	case []byte:
		pv.SetString(fmt.Sprintf("%s", value))
		return
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}

// Set obj to value's enum in stringMapEnum representation
func SQLScanEnumString(obj interface{}, value interface{}, stringMapEnum map[string]interface{}) (err error) {
	if value == nil {
		return
	}

	pv, err := ensureScanValue(obj)
	if err != nil {
		return
	}

	switch value.(type) {
	case []byte:
		enumstr := fmt.Sprintf("%s", value)
		enumval, ok := stringMapEnum[enumstr]
		if !ok {
			return ErrorNoValueFound1.New(nil, value)
		}
		ev := reflect.ValueOf(enumval)
		pv.Set(ev)
		return
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}

// Set obj to value's stringslice representation
func SQLScanStringSlice(obj interface{}, value interface{}) (err error) {
	if value == nil {
		return
	}

	pv, err := ensureScanValue(obj)
	if err != nil {
		return
	}

	switch value.(type) {
	case []byte:
		stringslice := parseArray(string(value.([]byte)))
		pv.Set(reflect.ValueOf(stringslice))
		return
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}

// Return s postgresql representation which implements driver.Value
func SQLValueStringSlice(obj interface{}) (value driver.Value, err error) {
	if obj == nil {
		return
	}

	rv := reflect.ValueOf(obj)

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return
		} else {
			return SQLValueStringSlice(reflect.Indirect(rv).Interface())
		}
	case reflect.Slice:
		switch rv.Type().Elem().Kind() {
		case reflect.String:
			strs := []string{}
			for i, rvlen := 0, rv.Len(); i < rvlen; i++ {
				rvelem := rv.Index(i)
				strelem := `"` + strings.Replace(strings.Replace(rvelem.String(), `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
				strs = append(strs, strelem)
			}
			return "{" + strings.Join(strs, ",") + "}", nil
		}
	}

	return "", ErrorInvalidObjType1.New(nil, obj)
}
