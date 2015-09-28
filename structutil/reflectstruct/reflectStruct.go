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

func ReflectStruct(dest interface{}, src interface{}) (err error) {
	if dest == nil || src == nil {
		return
	}

	destReflectMap_TagName_FieldIndex := map[string]int{}

	valdest := reflectutil.EnsureValue(reflect.ValueOf(dest))
	for i := 0; i < valdest.NumField(); i++ {
		tfield := valdest.Type().Field(i)
		tag := tfield.Tag
		tagjsons := strings.Split(tag.Get("json"), ",")
		if len(tagjsons) > 0 && tagjsons[0] != "-" {
			destReflectMap_TagName_FieldIndex[tagjsons[0]] = i
		}
	}

	valsrc := reflect.ValueOf(src)

	switch valsrc.Kind() {
	case reflect.Map:
		for _, valsrcmapkey := range valsrc.MapKeys() {
			switch valsrcmapkey.Kind() {
			case reflect.String:
				srcmapkey := valsrcmapkey.Interface().(string)
				if i, ok := destReflectMap_TagName_FieldIndex[srcmapkey]; ok {
					valsrcmapval := reflectutil.EnsureValue(valsrc.MapIndex(valsrcmapkey))
					valDestField := valdest.Field(i)
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
