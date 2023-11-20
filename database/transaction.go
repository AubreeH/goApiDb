package database

import "database/sql"

type Transaction struct {
	tx *sql.Tx
	db *Database
}
