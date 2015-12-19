package apiutil

import (
	"database/sql/driver"

	"github.com/tsaikd/KDGoLib/enumutil"
)

type Method int8

const (
	MethodAny Method = 1 + iota
	MethodGET
	MethodPOST
	MethodDELETE
	MethodPATCH
	MethodPUT
	MethodOPTIONS
	MethodHEAD
)

var methodEnum = enumutil.NewEnumFactory().
	Add(MethodAny, "Any").
	Add(MethodGET, "GET").
	Add(MethodPOST, "POST").
	Add(MethodDELETE, "DELETE").
	Add(MethodPATCH, "PATCH").
	Add(MethodPUT, "PUT").
	Add(MethodOPTIONS, "OPTIONS").
	Add(MethodHEAD, "HEAD").
	Build()

func (t Method) String() string {
	return methodEnum.String(t)
}

func (t Method) MarshalJSON() ([]byte, error) {
	return methodEnum.MarshalJSON(t)
}

func (t *Method) UnmarshalJSON(b []byte) (err error) {
	return methodEnum.UnmarshalJSON(t, b)
}

func (t *Method) Scan(value interface{}) (err error) {
	return methodEnum.Scan(t, value)
}

func (t Method) Value() (v driver.Value, err error) {
	return methodEnum.Value(t)
}

func IsMethod(s string) bool {
	return methodEnum.IsEnumString(s)
}

func ParseMethod(s string) Method {
	enum, err := methodEnum.ParseString(s)
	if err != nil {
		return 0
	}
	return enum.(Method)
}
