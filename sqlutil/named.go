package sqlutil

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
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

// NamedExecContext prepare named SQL statement and call Exec function with context
func NamedExecContext(
	ctx context.Context,
	tx dbtypes.Named,
	query string,
	arg interface{},
) (result sql.Result, err error) {
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return result, ErrorPrepareNamed.New(err)
	}
	if result, err = stmt.ExecContext(ctx, arg); err != nil {
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

// NamedExecStrictContext like NamedExecContext, but result rows affected should greater than zero
func NamedExecStrictContext(
	ctx context.Context,
	tx dbtypes.Named,
	query string,
	arg interface{},
) (result sql.Result, err error) {
	if result, err = NamedExecContext(ctx, tx, query, arg); err != nil {
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

// NamedGetContext prepare named SQL statement and call Get function with context
func NamedGetContext(
	ctx context.Context,
	tx dbtypes.Named,
	dest interface{},
	query string,
	arg interface{},
) (err error) {
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return ErrorPrepareNamed.New(err)
	}
	if err = stmt.GetContext(ctx, dest, arg); err != nil {
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

// NamedSelectContext prepare named SQL statement and call Select function with context
func NamedSelectContext(
	ctx context.Context,
	tx dbtypes.Named,
	dest interface{},
	query string,
	arg interface{},
) (err error) {
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return ErrorPrepareNamed.New(err)
	}
	if err = stmt.SelectContext(ctx, dest, arg); err != nil {
		return wrapError(err)
	}
	return nil
}

// NamedSelectInContext prepare named/in SQL statement and call Select function with context
func NamedSelectInContext(
	ctx context.Context,
	tx dbtypes.Transactionx,
	dest interface{},
	query string,
	arg interface{},
) (err error) {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return ErrorPrepareNamed.New(err)
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return ErrorPrepareNamed.New(err)
	}

	query = tx.Rebind(query)
	if err = tx.SelectContext(ctx, dest, query, args...); err != nil {
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
