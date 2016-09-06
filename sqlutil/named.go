package sqlutil

import (
	"database/sql"

	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/sqlutil/dbtypes"
)

// errors
var (
	ErrorPrepareNamed  = errutil.NewFactory("prepare SQL named query failed")
	ErrorNoRowAffected = errutil.NewFactory("no row affected")
)

// NamedExec prepare named SQL statement and call Exec function
func NamedExec(
	tx dbtypes.Named,
	query string,
	arg interface{},
) (result sql.Result, err error) {
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return result, ErrorPrepareNamed.New(err)
	}
	if result, err = stmt.Exec(arg); err != nil {
		return result, wrapError(err)
	}
	return result, nil
}

// NamedExecStrict like NamedExec, but result rows affected should greater than zero
func NamedExecStrict(
	tx dbtypes.Named,
	query string,
	arg interface{},
) (result sql.Result, err error) {
	if result, err = NamedExec(tx, query, arg); err != nil {
		return result, err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return result, err
	}
	if num < 1 {
		return result, ErrorNoRowAffected.New(nil)
	}
	return result, nil
}

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
		return wrapError(err)
	}
	return nil
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
		return wrapError(err)
	}
	return nil
}

func wrapError(err error) error {
	if IsErrorNoRowsInResultSet(err) {
		return ErrorNoRowsInResultSet.New(err)
	}
	if IsErrorDuplicateViolateUniqueConstraint(err) {
		return ErrorDuplicateViolateUniqueConstraint.New(err)
	}
	if IsErrorTsquerySyntax(err) {
		return ErrorTsquerySyntax.New(err)
	}
	return err
}
