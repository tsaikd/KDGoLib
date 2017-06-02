package reflectstruct

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/jsonex"
	"github.com/tsaikd/KDGoLib/reflectutil"
	"github.com/tsaikd/govalidator"
)

// expose errors
var (
	ErrorUnknownSourceMapKeyType1       = errutil.NewFactory("unknown source map key type: %v")
	ErrorUnsupportedReflectFieldMethod2 = errutil.NewFactory("unsupported reflect field method: %v <- %v")
	ErrorFieldCanNotSet                 = errutil.NewFactory("field can not set")
)

var (
	jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	typeOfBytes     = reflect.TypeOf([]byte(nil))
)

func reflectField(field reflect.Value, val reflect.Value) (err error) {
	if field.IsValid() && field.Type().Implements(jsonUnmarshaler) {
		switch val.Kind() {
		case reflect.String:
			return json.Unmarshal([]byte(`"`+val.String()+`"`), field.Interface())
		case reflect.Slice:
			if val.Type() == typeOfBytes {
				return json.Unmarshal(val.Bytes(), field.Interface())
			}
		}
	}

	if field.CanAddr() && field.Addr().Type().Implements(jsonUnmarshaler) {
		switch val.Kind() {
		case reflect.String:
			return reflectField(field.Addr(), val)
		case reflect.Slice:
			if val.Type() == typeOfBytes {
				return reflectField(field.Addr(), val)
			}
		}
	}

	// get val real type
	valElemType := val.Type()
	switch valElemType.Kind() {
	case reflect.Ptr:
		valElemType = valElemType.Elem()
		val = val.Elem()
	case reflect.Interface:
		val = val.Elem()
	}

	// field is not slice, val is slice, assign last val
	if field.Kind() != reflect.Slice && val.Kind() == reflect.Slice {
		lastidx := val.Len() - 1
		if lastidx >= 0 {
			return reflectField(field, val.Index(lastidx))
		}
		return
	}

	switch field.Kind() {
	case reflect.Ptr:
		if field.IsNil() {
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.Set(reflect.New(field.Type().Elem()))
		}
		return reflectField(field.Elem(), val)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		switch val.Kind() {
		case reflect.String:
			valstr := val.String()
			var v int64
			if v, err = strconv.ParseInt(valstr, 0, 64); err != nil {
				return
			}
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetInt(v)
			return
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetInt(val.Int())
			return
		case reflect.Float64, reflect.Float32:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetInt(int64(val.Float()))
			return
		case reflect.Bool:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			if val.Bool() {
				field.SetInt(1)
			} else {
				field.SetInt(0)
			}
			return
		}
	case reflect.Float64, reflect.Float32:
		switch val.Kind() {
		case reflect.String:
			valstr := val.String()
			var v float64
			if v, err = strconv.ParseFloat(valstr, 64); err != nil {
				return
			}
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetFloat(v)
			return
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetFloat(float64(val.Int()))
			return
		case reflect.Float64, reflect.Float32:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetFloat(val.Float())
			return
		case reflect.Bool:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			if val.Bool() {
				field.SetFloat(1)
			} else {
				field.SetFloat(0)
			}
			return
		}
	case reflect.Bool:
		switch val.Kind() {
		case reflect.String:
			valstr := val.String()
			var v bool
			if v, err = govalidator.ToBoolean(valstr); err != nil {
				return
			}
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetBool(v)
			return
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetBool(val.Int() > 0)
			return
		case reflect.Float64, reflect.Float32:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetBool(val.Float() > 0)
			return
		case reflect.Bool:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetBool(val.Bool())
			return
		}
	case reflect.String:
		switch val.Kind() {
		case reflect.String:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetString(val.String())
			return
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetString(strconv.FormatInt(val.Int(), 10))
			return
		case reflect.Float64, reflect.Float32:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetString(strconv.FormatFloat(val.Float(), 'f', -1, 64))
			return
		case reflect.Bool:
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.SetString(strconv.FormatBool(val.Bool()))
			return
		}
	case reflect.Slice:
		switch val.Kind() {
		case reflect.Slice:
			len := val.Len()
			vals := reflect.MakeSlice(field.Type(), len, len)
			for i := 0; i < len; i++ {
				if err = reflectField(vals.Index(i), val.Index(i)); err != nil {
					return
				}
			}
			if !field.CanSet() {
				return ErrorFieldCanNotSet.New(nil)
			}
			field.Set(vals)
			return
		}
	case reflect.Struct:
		fieldInfoMap := buildReflectFieldInfo(nil, field)

		switch val.Kind() {
		case reflect.Map:
			for _, valmapkey := range val.MapKeys() {
				switch valmapkey.Kind() {
				case reflect.String:
					mapkey := valmapkey.String()
					if valDestField, ok := fieldInfoMap[mapkey]; ok {
						valmapval := reflectutil.EnsureValue(val.MapIndex(valmapkey))
						if err = reflectField(valDestField, valmapval); err != nil {
							return
						}
					}
				default:
					return ErrorUnknownSourceMapKeyType1.New(nil, valmapkey.Kind())
				}
			}
			return
		case reflect.Struct:
			if val.Type().AssignableTo(field.Type()) {
				if !field.CanSet() {
					return ErrorFieldCanNotSet.New(nil)
				}
				field.Set(val)
				return
			}

			valInfoMap := buildReflectFieldInfo(nil, val)
			for key, valField := range valInfoMap {
				fieldInfo, exist := fieldInfoMap[key]
				if !exist {
					continue
				}
				if err = reflectField(fieldInfo, valField); err != nil {
					return
				}
			}
			return
		case reflect.String:
			if field.CanAddr() && field.Addr().CanInterface() && field.Addr().Interface() != nil {
				return jsonex.Unmarshal([]byte(val.String()), field.Addr().Interface())
			}
		}
	case reflect.Interface:
		if !field.CanSet() {
			return ErrorFieldCanNotSet.New(nil)
		}
		field.Set(val)
		return
	case reflect.Map:
		switch val.Kind() {
		case reflect.Map:
			if !val.IsNil() {
				if !field.CanSet() {
					return ErrorFieldCanNotSet.New(nil)
				}
				field.Set(val)
			}
			return
		case reflect.Struct:
			return reflectFieldStruct2Map(field, val)
		}
	}

	return ErrorUnsupportedReflectFieldMethod2.New(nil, field.Kind(), val.Kind())
}

func buildReflectFieldInfo(fieldInfoMap map[string]reflect.Value, value reflect.Value) map[string]reflect.Value {
	if fieldInfoMap == nil {
		fieldInfoMap = map[string]reflect.Value{}
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if tag := field.Tag.Get("json"); tag != "" {
			tagvals := strings.Split(tag, ",")
			if len(tagvals) > 0 {
				if tagvals[0] != "-" {
					fieldInfoMap[tagvals[0]] = value.Field(i)
				}
				continue
			}
		}
		if tag := field.Tag.Get("reflect"); tag == "inherit" {
			childValue := value.Field(i)
			if childValue.Kind() == reflect.Ptr && childValue.IsNil() {
				childValue = reflect.New(childValue.Type().Elem())
				if value.Field(i).CanSet() {
					value.Field(i).Set(childValue)
				}
				childValue = childValue.Elem()
			}
			fieldInfoMap = buildReflectFieldInfo(fieldInfoMap, childValue)
			continue
		}
		if field.Name[:1] == strings.ToUpper(field.Name[:1]) {
			fieldInfoMap[field.Name] = value.Field(i)
			continue
		}
	}
	return fieldInfoMap
}

// ReflectStruct reflect data from src to dest
func ReflectStruct(dest interface{}, src interface{}) (err error) {
	if dest == nil || src == nil {
		return
	}

	return reflectField(
		reflect.ValueOf(dest),
		reflect.ValueOf(src),
	)
}
