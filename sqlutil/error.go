package sqlutil

import (
	"database/sql"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorNoRowsInResultSet                = errutil.NewFactory("no rows in result set")
	ErrorDuplicateViolateUniqueConstraint = errutil.NewFactory("duplicate key value violates unique constraint")
	ErrorTsquerySyntax                    = errutil.NewFactory("syntax error in tsquery")
)

// IsErrorNoRowsInResultSet check err is sql error "no rows in result set" or "no row affected"
func IsErrorNoRowsInResultSet(err error) bool {
	switch err {
	case nil:
		return false
	case sql.ErrNoRows:
		return true
	}
	if ErrorNoRowsInResultSet.Match(err) {
		return true
	}
	switch err.Error() {
	case "sql: no rows in result set", "no row affected":
		return true
	}
	return false
}

// IsContainErrorNoRowsInResultSet check err contain sql error "no rows in result set"
func IsContainErrorNoRowsInResultSet(err error) bool {
	return errutil.ContainErrorFunc(err, IsErrorNoRowsInResultSet)
}

// IsErrorDuplicateViolateUniqueConstraint check err is sql error "duplicate key value violates unique constraint"
func IsErrorDuplicateViolateUniqueConstraint(err error) bool {
	if ErrorDuplicateViolateUniqueConstraint.Match(err) {
		return true
	}
	switch e := err.(type) {
	case nil:
		return false
	case *pq.Error:
		return e.Code == "23505"
	case *mysql.MySQLError:
		return e.Number == 1062
	default:
		return strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint")
	}
}

// IsContainErrorDuplicateViolateUniqueConstraint check err contain sql error "duplicate key value violates unique constraint"
func IsContainErrorDuplicateViolateUniqueConstraint(err error) bool {
	return errutil.ContainErrorFunc(err, IsErrorDuplicateViolateUniqueConstraint)
}

// IsErrorTsquerySyntax check error is sql error "syntax error in tsquery"
func IsErrorTsquerySyntax(err error) bool {
	if ErrorTsquerySyntax.Match(err) {
		return true
	}
	switch err.(type) {
	case nil:
		return false
	case *pq.Error:
		e := err.(*pq.Error)
		if e.Code != "42601" {
			return false
		}
		return strings.Contains(e.Message, "syntax error in tsquery:")
	default:
		return strings.Contains(err.Error(), "pq: syntax error in tsquery:")
	}
}

// IsContainErrorTsquerySyntax check error contain sql error "syntax error in tsquery"
func IsContainErrorTsquerySyntax(err error) bool {
	return errutil.ContainErrorFunc(err, IsErrorTsquerySyntax)
}
