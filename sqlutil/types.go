package sqlutil

type SQLStringSlice []string

func (t *SQLStringSlice) Scan(value interface{}) (err error) {
	return SQLScanStringSlice(t, value)
}
