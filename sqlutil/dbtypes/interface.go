package dbtypes

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/jmoiron/sqlx"
)

// Arg for stmt argument
type Arg interface {
	Value() (value driver.Value, err error)
}

// Result for select column
type Result interface {
	Scan(value interface{}) (err error)
}

// Exec sql stmt
type Exec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Get is select only one row
type Get interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Select multiple rows
type Select interface {
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Named support named argument in stmt
type Named interface {
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// Queryx support query row
type Queryx interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

// Rebind support rebind args
type Rebind interface {
	Rebind(query string) string
}

// Transaction for db connection
type Transaction interface {
	Commit() error
	Rollback() error

	Exec
	Get
	Select
	Named
}

// Transactionx more interface for Transaction
type Transactionx interface {
	Queryx
	Rebind

	Transaction
}

// DBLike is general db object
type DBLike interface {
	Exec
	Get
	Select
	Named
}
