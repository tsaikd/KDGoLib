package sqlutil

import "github.com/jmoiron/sqlx"

// ToTsQuery return DB accepted query text
func ToTsQuery(db *sqlx.DB, query string) (tsquery string, err error) {
	if err = db.Get(&tsquery, `SELECT to_tsquery($1);`, query); IsErrorTsquerySyntax(err) {
		err = db.Get(&tsquery, `SELECT plainto_tsquery($1);`, query)
	}
	if tsquery == "" && query != "" {
		tsquery = query
	}
	return
}
