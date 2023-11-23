package driver

import "database/sql"

func (d *DriverType) DescribeTable(db *sql.DB, tableName string) (*TableDescription, error) {
	switch *d {
	case MySql:
		return describeMySqlTable(db, tableName)
	case Postgres:
		return describePostgresTable(db, tableName)
	case SQLite:
		return describeSQLiteTable(db, tableName)
	default:
		return nil, nil
	}
}
