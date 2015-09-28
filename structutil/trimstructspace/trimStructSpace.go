package trimstructspace

import (
	"reflect"
	"strings"

	"github.com/tsaikd/KDGoLib/reflectutil"
)

func TrimStructSpace(dest interface{}) (err error) {
	if dest == nil {
		return
	}

	valdest := reflectutil.EnsureValue(reflect.ValueOf(dest))
	for i := 0; i < valdest.NumField(); i++ {
		field := valdest.Field(i)
		switch field.Kind() {
		case reflect.String:
			if field.CanSet() {
				str := field.String()
				str = strings.TrimSpace(str)
				field.SetString(str)
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				len := field.Len()
				vals := reflect.MakeSlice(field.Type(), len, len)
				for i := 0; i < len; i++ {
					str := field.Index(i).Interface().(string)
					str = strings.TrimSpace(str)
					vals.Index(i).SetString(str)
				}
				field.Set(vals)
			}
		case reflect.Map:
			for _, valkey := range field.MapKeys() {
				valval := field.MapIndex(valkey)
				if valval.Kind() == reflect.String {
					str := valval.Interface().(string)
					str = strings.TrimSpace(str)
					field.SetMapIndex(valkey, reflect.ValueOf(str))
				}
			}
		case reflect.Struct:
			if err = TrimStructSpace(field.Addr().Interface()); err != nil {
				return
			}
		}
	}

	return
}
