package sqlutil

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tsaikd/KDGoLib/errutil"
)

var (
	ErrorUnsupportedScanType1 = errutil.ErrorFactory("unsupported scan type: %v")
	ErrorInvalidObjType1      = errutil.ErrorFactory("invalid obj type: %T")
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
	if value, err = json.Marshal(obj); err != nil {
		return
	}
	return
}

// Set obj to value's string representation
func SQLScanString(obj interface{}, value interface{}) (err error) {
	if value == nil {
		return
	}

	// Get the value of obj and make sure it's either a pointer or nil
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrorInvalidObjType1.New(nil, obj)
	}
	// So we can get its actual value
	pv := reflect.Indirect(rv)

	switch value.(type) {
	case []byte:
		pv.SetString(fmt.Sprintf("%s", value))
		return
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}
