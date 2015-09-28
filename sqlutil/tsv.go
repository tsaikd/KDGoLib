package sqlutil

import "github.com/jmoiron/sqlx"

func ToTsQuery(db *sqlx.DB, query string) (tsquery string, err error) {
	err = db.Get(
		&tsquery,
		`	SELECT to_tsquery($1)
		;`,
		query,
	)
	if IsErrorTsquerySyntax(err) {
		err = db.Get(
			&tsquery,
			`	SELECT plainto_tsquery($1)
			;`,
			query,
		)
	}
	if tsquery == "" && query != "" {
		tsquery = query
	}
	return
}
