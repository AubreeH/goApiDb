package access

import (
	"fmt"
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/structParsing"
)

func Create[T any](db *database.Database, values []T) error {
	var entity T
	tableInfo, err := structParsing.GetTableInfo(entity)

	queryColumns := ""
	queryValues := ""
	var args []any

	for i := range values {
		var rowData []ColumnData
		rowData, err = GetData(values[i], createOperationHandler)
		if err != nil {
			return err
		}

		for j := range rowData {
			columnData := rowData[j]

			if queryColumns == "" {
				queryColumns += columnData.Column
				queryValues += "?"
			} else {
				queryColumns += ", " + columnData.Column
				queryValues += ", ?"
			}

			args = append(args, columnData.Data)
		}
	}

	query := fmt.Sprintf("INSERT"+" INTO %s (%s) values (%s)", tableInfo.Name, queryColumns, queryValues)

	_, err = db.Db.Exec(query, args...)
	if err != nil {
		return err
	}

	return err
}
