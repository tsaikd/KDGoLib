package sqlutil

import (
	"database/sql/driver"

	"github.com/lib/pq"
	"github.com/tsaikd/KDGoLib/jsonex"
)

// SQLStringSlice is a type of []string and implement SQL driver
type SQLStringSlice []string

// Scan decode SQL value
func (t *SQLStringSlice) Scan(value interface{}) (err error) {
	return SQLScanStringSlice(t, value)
}

// Value return data for SQL
func (t *SQLStringSlice) Value() (value driver.Value, err error) {
	return SQLValueStringSlice(t)
}

// ToStringSlice return nil if SQLStringSlice is nil
func (t *SQLStringSlice) ToStringSlice() *[]string {
	if t == nil {
		return nil
	}
	data := []string(*t)
	return &data
}

// NewSQLStringSlice create from *[]string
func NewSQLStringSlice(s *[]string) *SQLStringSlice {
	if s == nil {
		return nil
	}
	data := SQLStringSlice(*s)
	return &data
}

// SQLStringSliceJSON is a type of []string and implement SQL driver
type SQLStringSliceJSON []string

// Scan decode SQL value
func (t *SQLStringSliceJSON) Scan(value interface{}) (err error) {
	return SQLScanJSON(t, value)
}

// Value return data for SQL
func (t *SQLStringSliceJSON) Value() (value driver.Value, err error) {
	return SQLValueJSON(t)
}

// SQLJsonMap general map type for SQL
type SQLJsonMap map[string]interface{}

// Scan decode SQL value
func (t *SQLJsonMap) Scan(value interface{}) (err error) {
	return SQLScanJSON(t, value)
}

// Value return data for SQL
func (t *SQLJsonMap) Value() (value driver.Value, err error) {
	return SQLValueJSON(t)
}

// SQLNullTime represents a time.Time that may be null. Inherit pq.NullTime and rewrite JSON marshaler
type SQLNullTime struct {
	pq.NullTime
}

// MarshalJSON implement JSON marshaler
func (t SQLNullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return jsonex.Marshal(nil)
	}
	return jsonex.Marshal(t.Time)
}

// UnmarshalJSON implement JSON unmarshaler
func (t *SQLNullTime) UnmarshalJSON(data []byte) (err error) {
	if err = t.Time.UnmarshalJSON(data); err != nil {
		t.Valid = false
		return
	}
	t.Valid = true
	return nil
}
