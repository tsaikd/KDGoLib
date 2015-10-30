package sqlutil

import "database/sql/driver"

type SQLStringSlice []string

func (t *SQLStringSlice) Scan(value interface{}) (err error) {
	return SQLScanStringSlice(t, value)
}

func (t *SQLStringSlice) Value() (value driver.Value, err error) {
	return SQLValueStringSlice(t)
}

type SQLJsonMap map[string]interface{}

func (t *SQLJsonMap) Scan(value interface{}) (err error) {
	return SQLScanJson(t, value)
}

func (t *SQLJsonMap) Value() (value driver.Value, err error) {
	return SQLValueJson(t)
}
