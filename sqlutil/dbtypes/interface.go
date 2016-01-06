package dbtypes

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Exec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Get interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

type Select interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

type Named interface {
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type Queryx interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

type Transaction interface {
	Commit() error
	Rollback() error

	Exec
	Get
	Select
	Named
}

type Transactionx interface {
	Queryx

	Transaction
}

type DBLike interface {
	Exec
	Get
	Select
	Named
}
