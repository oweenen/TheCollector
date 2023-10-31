package database

import "database/sql"

func NewTransaction() (*sql.Tx, error) {
	return db.Begin()
}
