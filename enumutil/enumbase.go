package enumutil

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/sqlutil"
)

// errors
var (
	ErrorEnumString1 = errutil.NewFactory("convert string to enum failed: %v")
)

type EnumString interface {
	String() string
}

type enumBase struct {
	mape2s map[EnumString]string
	maps2e map[string]interface{}
}

func (t enumBase) String(enum EnumString) string {
	return t.mape2s[enum]
}

func (t enumBase) ParseString(str string) (interface{}, error) {
	enum, ok := t.maps2e[str]
	if !ok {
		return 0, ErrorEnumString1.New(nil, str)
	}
	return enum, nil
}

func (t enumBase) MarshalJSON(enum EnumString) ([]byte, error) {
	return []byte(`"` + enum.String() + `"`), nil
}

func (t enumBase) UnmarshalJSON(enum interface{}, b []byte) (err error) {
	str := ""
	if err = json.Unmarshal(b, &str); err != nil {
		return
	}
	en, ok := t.maps2e[str]
	if !ok {
		return errutil.New("can not unmarshal data to enum, " + string(b))
	}

	reflect.ValueOf(enum).Elem().Set(reflect.ValueOf(en))
	return nil
}

func (t enumBase) Scan(enum interface{}, value interface{}) (err error) {
	return sqlutil.SQLScanEnumString(enum, value, t.maps2e)
}

func (t enumBase) Value(enum EnumString) (v driver.Value, err error) {
	v = enum.String()
	return
}

func (t enumBase) IsEnumString(str string) bool {
	_, ok := t.maps2e[str]
	return ok
}

func (t enumBase) Each(walker func(enum interface{}) (stop bool, err error)) error {
	for _, enum := range t.maps2e {
		stop, err := walker(enum)
		if stop || err != nil {
			return err
		}
	}
	return nil
}
