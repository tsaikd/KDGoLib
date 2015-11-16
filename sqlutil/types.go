package sqlutil

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type SQLStringSlice []string

func (t *SQLStringSlice) Scan(value interface{}) (err error) {
	return SQLScanStringSlice(t, value)
}

func (t *SQLStringSlice) Value() (value driver.Value, err error) {
	return SQLValueStringSlice(t)
}

func (t *SQLStringSlice) ToStringSlice() *[]string {
	if t == nil {
		return nil
	}
	data := []string(*t)
	return &data
}

func NewSQLStringSlice(s *[]string) *SQLStringSlice {
	if s == nil {
		return nil
	}
	data := SQLStringSlice(*s)
	return &data
}

type SQLJsonMap map[string]interface{}

func (t *SQLJsonMap) Scan(value interface{}) (err error) {
	return SQLScanJson(t, value)
}

func (t *SQLJsonMap) Value() (value driver.Value, err error) {
	return SQLValueJson(t)
}

// JsonText is a json.RawMessage, which is a []byte underneath.
// Value() validates the json format in the source, and returns an error if
// the json is not valid.  Scan does no validation.  JsonText additionally
// implements `Unmarshal`, which unmarshals the json within to an interface{}
type JsonText json.RawMessage

// MarshalJSON returns the *j as the JSON encoding of j.
func (j *JsonText) MarshalJSON() ([]byte, error) {
	return *j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JsonText) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JsonText: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil

}

// Value returns j as a value.  This does a validating unmarshal into another
// RawMessage.  If j is invalid json, it returns an error.
func (j *JsonText) Value() (value driver.Value, err error) {
	if j == nil {
		return
	}
	var m json.RawMessage
	err = j.Unmarshal(&m)
	if err != nil {
		return []byte{}, err
	}
	return []byte(*j), nil
}

// Scan stores the src in *j.  No validation is done.
func (j *JsonText) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Incompatible type for JsonText")
	}
	*j = JsonText(append((*j)[0:0], source...))
	return nil
}

// Unmarshal unmarshal's the json in j to v, as in json.Unmarshal.
func (j *JsonText) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(*j), v)
}
