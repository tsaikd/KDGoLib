package sqlutil

import (
	"strings"

	"github.com/lib/pq"
	"github.com/tsaikd/KDGoLib/errutil"
)

func IsErrorNoRowsInResultSet(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "sql: no rows in result set"
}

func IsContainErrorNoRowsInResultSet(err error) bool {
	return errutil.ContainErrorFunc(err, IsErrorNoRowsInResultSet)
}

func IsErrorDuplicateViolateUniqueConstraint(err error) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
	case *pq.Error:
		e := err.(*pq.Error)
		if e.Code != "23505" {
			return false
		}
		return strings.Contains(e.Message, "duplicate key value violates unique constraint")
	default:
		return strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint")
	}
}

func IsContainErrorDuplicateViolateUniqueConstraint(err error) bool {
	return errutil.ContainErrorFunc(err, IsErrorDuplicateViolateUniqueConstraint)
}

func IsErrorTsquerySyntax(err error) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
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

func IsContainErrorTsquerySyntax(err error) bool {
	return errutil.ContainErrorFunc(err, IsErrorTsquerySyntax)
}
