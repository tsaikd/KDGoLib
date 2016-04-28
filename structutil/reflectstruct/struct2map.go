package reflectstruct

import "reflect"

func reflectFieldStruct2Map(field reflect.Value, val reflect.Value) (err error) {
	valInfoMap := buildReflectFieldInfo(nil, val)
	for key, valField := range valInfoMap {
		field.SetMapIndex(reflect.ValueOf(key), valField)
	}
	return
}
