package sqlutil

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// ParseTsQuery remove invalid tsquery syntax, tx should not be transaction
// because error will cause transaction stop
func ParseTsQuery(ctx context.Context, tx *sqlx.DB, query string) (result string, err error) {
	if err = tx.GetContext(ctx, &result, `SELECT to_tsquery($1);`, query); err == nil {
		return result, nil
	}

	// fallback with plainto_tsquery()
	if err = tx.GetContext(ctx, &result, `SELECT plainto_tsquery($1);`, query); err == nil {
		return result, nil
	}

	return query, err
}
