package driver

import "database/sql"

func (d *DriverType) GetKey(db *sql.DB, tableName, columnName string) (*Key, error) {
	switch *d {
	case MySql:
		return getMySqlKey(db, tableName, columnName)
	case Postgres:
		return getPostgresKey(db, tableName, columnName)
	case SQLite:
		return getSQLiteKey(db, tableName, columnName)
	default:
		return nil, nil
	}
}
