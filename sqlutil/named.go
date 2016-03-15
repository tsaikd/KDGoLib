package sqlutil

import (
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/sqlutil/dbtypes"
)

// errors
var (
	ErrorPrepareNamed = errutil.NewFactory("prepare SQL named query failed")
)

// NamedGet prepare named SQL statement and call Get function
func NamedGet(
	tx dbtypes.Named,
	dest interface{},
	query string,
	arg interface{},
) (err error) {
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return ErrorPrepareNamed.New(err)
	}
	if err = stmt.Get(dest, arg); err != nil {
		return
	}
	return
}

// NamedSelect prepare named SQL statement and call Select function
func NamedSelect(
	tx dbtypes.Named,
	dest interface{},
	query string,
	arg interface{},
) (err error) {
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return ErrorPrepareNamed.New(err)
	}
	if err = stmt.Select(dest, arg); err != nil {
		return
	}
	return
}
