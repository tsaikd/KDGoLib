package sqlutil

import (
	"strings"

	"github.com/lib/pq"
)

func IsErrorNoRowsInResultSet(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "sql: no rows in result set"
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
