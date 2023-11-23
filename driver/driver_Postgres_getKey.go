package driver

import (
	"database/sql"
	"errors"
)

func getPostgresKey(db *sql.DB, tableName string, columnName string) (*Key, error) {
	row := db.QueryRow(`
		SELECT
			TABLE_CATALOG,
			TABLE_SCHEMA,
			TABLE_NAME,
			COLUMN_NAME,
			CONSTRAINT_CATALOG,
			CONSTRAINT_SCHEMA,
			CONSTRAINT_NAME,
			ORDINAL_POSITION
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
		WHERE TABLE_NAME = ?
		AND COLUMN_NAME = ?
	`, tableName, columnName)

	var key Key
	err := row.Scan(
		&key.TableCatalog,
		&key.TableSchema,
		&key.TableName,
		&key.ColumnName,
		&key.ConstraintCatalog,
		&key.ConstraintSchema,
		&key.ConstraintName,
		&key.OrdinalPosition,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &key, nil
}
