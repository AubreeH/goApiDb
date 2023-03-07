package access

import (
	"fmt"
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/entities"
)

func Create[T any](db *database.Database, values []T) (T, error) {
	var entity T
	tableInfo, err := entities.GetTableInfo(entity)

	queryColumns := ""
	queryValues := ""
	var args []any

	var output T
	var id any
	for i := range values {
		var rowData []ColumnData
		rowData, err = GetData(values[i], CreateOperationHandler)
		if err != nil {
			return output, err
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

			if columnData.PrimaryKey {
				id = columnData.Data
			}

			args = append(args, columnData.Data)
		}
	}

	query := fmt.Sprintf("INSERT"+" INTO %s (%s) values (%s)", tableInfo.Name, queryColumns, queryValues)

	res, err := db.Db.Exec(query, args...)
	if err != nil {
		return output, err
	}

	if !(id == nil || id == 0 || id == "") {
		return GetById(db, output, id)
	}

	intId, err := res.LastInsertId()
	if err != nil {
		return output, err
	}

	newEntity, err := GetById(db, output, intId)

	return newEntity, err
}
