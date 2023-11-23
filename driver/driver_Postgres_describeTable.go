package driver

import "database/sql"

func describePostgresTable(db *sql.DB, tableName string) (*TableDescription, error) {
	rows, err := db.Query(`
		SELECT
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			tc.constraint_type
		FROM information_schema.columns c
		LEFT JOIN information_schema.key_column_usage kcu ON kcu.table_schema = c.table_schema AND kcu.table_name = c.table_name AND kcu.column_name = c.column_name
		LEFT JOIN information_schema.table_constraints tc ON tc.table_schema = kcu.table_schema AND tc.table_name = kcu.table_name AND tc.constraint_name = kcu.constraint_name
		WHERE table_name = $1
	`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var column Column
		err := rows.Scan(&column.Name, &column.Type, &column.Nullable, &column.DefaultValue, &column.Key)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return &TableDescription{
		Columns: columns,
	}, nil
}
