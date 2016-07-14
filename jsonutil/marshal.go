package jsonutil

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/tsaikd/KDGoLib/errutil"
)

var jsonNull, _ = json.Marshal(nil)

// errors
var (
	ErrorUnsupportedMarshalType2 = errutil.NewFactory("unsupported marshal type: %s (%s)")
)

// IsEmpty implement IsEmpty() instance
type IsEmpty interface {
	IsEmpty() bool
}

// MarshalJSON marshal JSON with extended flags
func MarshalJSON(v interface{}) ([]byte, error) {
	val, _, err := marshalValue(reflect.ValueOf(v), reflect.StructField{
		Tag: `json:",wrapper"`,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(val)
}

func marshalValue(refv reflect.Value, field reflect.StructField) (v interface{}, empty bool, err error) {
	if !refv.IsValid() {
		return nil, true, nil
	}

	_, tagmap := parseFieldJSONTag(field)

	if refv.CanInterface() {
		refval := refv.Interface()
		if reffunc, ok := refval.(IsEmpty); ok {
			if reffunc.IsEmpty() {
				return nil, true, nil
			}
		}
		if reffunc, ok := refval.(json.Marshaler); ok && !tagmap["wrapper"] {
			var buf []byte
			if buf, err = reffunc.MarshalJSON(); err != nil {
				return nil, true, err
			}
			if tagmap["omitempty"] && string(buf) == string(jsonNull) {
				return nil, true, nil
			}
			return buf, false, nil
		}
	}

	switch refv.Kind() {
	case reflect.Bool:
		buf := refv.Bool()
		v = buf
		def := false
		if tagmap["omitdefault"] {
			if def, err = strconv.ParseBool(field.Tag.Get("default")); err != nil {
				return nil, true, err
			}
			empty = buf == def
		} else if tagmap["omitempty"] {
			empty = buf == def
		} else {
			empty = false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf := refv.Int()
		v = buf
		def := int64(0)
		if tagmap["omitdefault"] {
			if def, err = strconv.ParseInt(field.Tag.Get("default"), 0, 64); err != nil {
				return nil, true, err
			}
			empty = buf == def
		} else if tagmap["omitempty"] {
			empty = buf == def
		} else {
			empty = false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		buf := refv.Uint()
		v = buf
		def := uint64(0)
		if tagmap["omitdefault"] {
			if def, err = strconv.ParseUint(field.Tag.Get("default"), 0, 64); err != nil {
				return nil, true, err
			}
			empty = buf == def
		} else if tagmap["omitempty"] {
			empty = buf == def
		} else {
			empty = false
		}
	case reflect.Float32, reflect.Float64:
		buf := refv.Float()
		v = buf
		def := float64(0)
		if tagmap["omitdefault"] {
			if def, err = strconv.ParseFloat(field.Tag.Get("default"), 64); err != nil {
				return nil, true, err
			}
			empty = buf == def
		} else if tagmap["omitempty"] {
			empty = buf == def
		} else {
			empty = false
		}
	case reflect.String:
		buf := refv.String()
		v = buf
		def := ""
		if tagmap["omitdefault"] {
			def = field.Tag.Get("default")
			empty = buf == def
		} else if tagmap["omitempty"] {
			empty = buf == def
		} else {
			empty = false
		}
	case reflect.Ptr, reflect.Interface:
		return marshalValue(refv.Elem(), field)
	case reflect.Slice:
		buf := []interface{}{}
		for i, n := 0, refv.Len(); i < n; i++ {
			elem, elemEmpty, elemErr := marshalValue(refv.Index(i), reflect.StructField{})
			if elemErr != nil {
				return nil, true, elemErr
			}
			if elemEmpty {
				continue
			}
			buf = append(buf, elem)
		}
		if len(buf) == 0 && tagmap["omitempty"] {
			v = nil
			empty = true
		} else {
			v = buf
			empty = false
		}
	case reflect.Struct:
		buf := map[string]interface{}{}
		reft := refv.Type()
		for i, max := 0, reft.NumField(); i < max; i++ {
			elemVal := refv.Field(i)
			elemField := reft.Field(i)
			if !elemField.Anonymous && !elemVal.CanInterface() {
				continue
			}
			elem, elemEmpty, elemErr := marshalValue(elemVal, elemField)
			if elemErr != nil {
				return nil, true, elemErr
			}
			if elemEmpty {
				continue
			}
			fieldname, _ := parseFieldJSONTag(elemField)
			if fieldname == "" {
				continue
			}
			if elemMap, ok := elem.(map[string]interface{}); ok && elemField.Anonymous {
				for elemk, elemv := range elemMap {
					buf[elemk] = elemv
				}
			} else {
				buf[fieldname] = elem
			}
		}
		if len(buf) == 0 && tagmap["omitempty"] {
			v = nil
			empty = true
		} else {
			v = buf
			empty = false
		}
	default:
		panic(ErrorUnsupportedMarshalType2.New(nil, refv.Kind(), refv.Type()))
	}

	return v, empty, nil
}

func parseFieldJSONTag(field reflect.StructField) (name string, tagmap map[string]bool) {
	tagmap = map[string]bool{}
	tag := field.Tag.Get("json")
	tags := strings.Split(tag, ",")

	if len(tags) < 1 {
		return field.Name, tagmap
	}

	name = strings.TrimSpace(tags[0])
	for i, n := 1, len(tags); i < n; i++ {
		tagmap[strings.TrimSpace(tags[i])] = true
	}

	switch name {
	case "-":
		return "", tagmap
	case "":
		return field.Name, tagmap
	default:
		return name, tagmap
	}
}
