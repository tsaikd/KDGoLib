package reflectstruct

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/reflectutil"
)

// expose errors
var (
	ErrorUnknownSourceMapKeyType1       = errutil.ErrorFactory("unknown source map key type: %v")
	ErrorUnknownSourceType1             = errutil.ErrorFactory("unknown source type: %v")
	ErrorUnsupportedReflectFieldMethod2 = errutil.ErrorFactory("unsupported reflect field method: %v <- %v")
)

func unsafeReflectFieldSlice2Slice(field reflect.Value, val reflect.Value) (err error) {
	if field.Type().Elem().Kind() == val.Type().Elem().Kind() {
		field.Set(val)
	} else {
		len := val.Len()
		vals := reflect.MakeSlice(field.Type(), len, len)
		for i := 0; i < len; i++ {
			if err = reflectField(vals.Index(i), val.Index(i)); err != nil {
				return
			}
		}
		field.Set(vals)
	}
	return
}

func reflectField(field reflect.Value, val reflect.Value) (err error) {
	val = reflectutil.EnsureValue(val)

	if field.Kind() == val.Kind() {
		switch field.Kind() {
		case reflect.Slice:
			return unsafeReflectFieldSlice2Slice(field, val)
		case reflect.String:
			field.SetString(val.String())
			return
		}

		field.Set(val)
		return
	}

	if val.Kind() == reflect.Slice {
		lastidx := val.Len() - 1
		if lastidx >= 0 {
			return reflectField(field, val.Index(lastidx))
		}
		return
	}

	switch field.Kind() {
	case reflect.Int64, reflect.Int:
		switch val.Kind() {
		case reflect.String:
			valstr := val.Interface().(string)
			v, err := strconv.ParseInt(valstr, 0, 64)
			if err != nil {
				return err
			}
			field.SetInt(v)
			return err
		case reflect.Int:
			field.SetInt(int64(val.Interface().(int)))
			return
		case reflect.Float64:
			field.SetInt(int64(val.Interface().(float64)))
			return
		}
	case reflect.Bool:
		switch val.Kind() {
		case reflect.String:
			valstr := val.Interface().(string)
			v, err := govalidator.ToBoolean(valstr)
			if err != nil {
				return err
			}
			field.SetBool(v)
			return nil
		}
	}

	return ErrorUnsupportedReflectFieldMethod2.New(nil, field.Kind(), val.Kind())
}

func buildReflectFieldInfo(fieldInfoMap map[string]reflect.Value, value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if tag := field.Tag.Get("json"); tag != "" {
			if tagvals := strings.Split(tag, ","); len(tagvals) > 0 && tagvals[0] != "-" {
				fieldInfoMap[tagvals[0]] = value.Field(i)
				continue
			}
		}
		if tag := field.Tag.Get("reflect"); tag == "inherit" {
			childValue := value.Field(i)
			if childValue.Kind() == reflect.Ptr && childValue.IsNil() {
				childValue = reflect.New(childValue.Type().Elem())
				value.Field(i).Set(childValue)
				childValue = childValue.Elem()
			}
			buildReflectFieldInfo(fieldInfoMap, childValue)
		}
	}
}

// ReflectStruct reflect data from src to dest
func ReflectStruct(dest interface{}, src interface{}) (err error) {
	if dest == nil || src == nil {
		return
	}

	fieldInfoMap := map[string]reflect.Value{}
	valdest := reflectutil.EnsureValue(reflect.ValueOf(dest))
	buildReflectFieldInfo(fieldInfoMap, valdest)

	valsrc := reflect.ValueOf(src)

	switch valsrc.Kind() {
	case reflect.Map:
		for _, valsrcmapkey := range valsrc.MapKeys() {
			switch valsrcmapkey.Kind() {
			case reflect.String:
				srcmapkey := valsrcmapkey.Interface().(string)
				if valDestField, ok := fieldInfoMap[srcmapkey]; ok {
					valsrcmapval := reflectutil.EnsureValue(valsrc.MapIndex(valsrcmapkey))
					if err = reflectField(valDestField, valsrcmapval); err != nil {
						return
					}
				}
			default:
				return ErrorUnknownSourceMapKeyType1.New(nil, valsrcmapkey.Kind())
			}
		}
	case reflect.Struct:
		data, err := json.Marshal(valsrc.Interface())
		if err != nil {
			return err
		}
		newstruct := map[string]interface{}{}
		if err = json.Unmarshal(data, &newstruct); err != nil {
			return err
		}
		return ReflectStruct(dest, newstruct)
	default:
		return ErrorUnknownSourceType1.New(nil, valsrc.Kind())
	}

	return
}
