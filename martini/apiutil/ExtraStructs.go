package apiutil

import "reflect"

// ExtraStruct contain struct name, type info, and params info
type ExtraStruct struct {
	FieldName string
	FieldType reflect.Type
	Params    RequestParams
}

// ExtraStructs is slice type of ExtraStruct
type ExtraStructs []*ExtraStruct

// Upsert append ExtraStruct to ExtraStructs if not exist
func (p *ExtraStructs) Upsert(param string, paramType reflect.Type) *ExtraStruct {
	t := *p
	for i, extraStruct := range t {
		if extraStruct.FieldName == param {
			t[i].FieldType = paramType
			return t[i]
		}
	}
	extraStruct := &ExtraStruct{
		FieldName: param,
	}
	if paramType != nil {
		extraStruct.FieldType = paramType
	}
	*p = append(t, extraStruct)
	return (*p)[len(*p)-1]
}
