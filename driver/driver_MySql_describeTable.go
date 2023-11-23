package driver

import (
	"database/sql"

	"github.com/AubreeH/goApiDb/structParsing"
)

func describeMySqlTable(db *sql.DB, tableName string) (*TableDescription, error) {
	rows, err := db.Query("DESCRIBE " + tableName)
	if err != nil {
		return nil, err
	}

	var table TableDescription
	for rows.Next() {
		var column Column
		err = rows.Scan(&column.Name, &column.Type, &column.Nullable, &column.Key, &column.DefaultValue, &column.Extra)
		if err != nil {
			return nil, err
		}

		column.Type = structParsing.FormatType(column.Type)
		column.DefaultValue = structParsing.FormatDefault(column.DefaultValue)
		column.Nullable = structParsing.FormatNullable(column.Nullable)
		column.Key = structParsing.FormatKey(column.Key)

		table.Columns = append(table.Columns, column)
	}

	return &table, nil
}
