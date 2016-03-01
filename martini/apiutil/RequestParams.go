package apiutil

import (
	"reflect"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorRequestParamNotFound1 = errutil.NewFactory("param %q not found in request map")
)

// RequestParam contain param name and type info
type RequestParam struct {
	FieldName string
	FieldType reflect.StructField
}

// RequestParams is slice type of RequestParam
type RequestParams []RequestParam

// Upsert append RequestParam to RequestParams if not exist
func (p *RequestParams) Upsert(param string, paramType *reflect.StructField) {
	if p == nil {
		return
	}
	t := *p
	for i, reqparam := range t {
		if reqparam.FieldName == param {
			if paramType != nil {
				t[i].FieldType = *paramType
			}
			return
		}
	}
	reqparam := RequestParam{
		FieldName: param,
	}
	if paramType != nil {
		reqparam.FieldType = *paramType
	}
	*p = append(t, reqparam)
}

func (p *RequestParams) delete(param string) (err error) {
	if p == nil {
		return
	}
	t := *p
	for i, reqparam := range t {
		if reqparam.FieldName == param {
			*p = append(t[:i], t[i+1:]...)
			return
		}
	}
	err = ErrorRequestParamNotFound1.New(nil, param)
	return
}

// Delete param in RequestParams
func (p *RequestParams) Delete(params ...string) (err error) {
	if p == nil {
		return
	}
	for _, param := range params {
		if err = p.delete(param); err != nil {
			return
		}
	}
	return
}

// Clone a new RequestParams
func (p *RequestParams) Clone() RequestParams {
	return *p
}

// FieldNames return field names in string slice type
func (p *RequestParams) FieldNames() (results []string) {
	if p == nil {
		return
	}
	t := *p
	for _, reqparam := range t {
		results = append(results, reqparam.FieldName)
	}
	return
}
