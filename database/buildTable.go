package database

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/helpers"
	"log"
	"strings"
)

func BuildTable(db *Database, entity interface{}, doLog bool, exec bool) error {
	tableName := helpers.GetTableName(entity)
	dbTableDescription, success, err := getTableDescription(db.Db, tableName)
	if err != nil {
		return err
	}

	structTableDescription := helpers.ConvertStructToTableDescription(entity)

	var rawSql string
	var constraints []string
	if success {
		rawSql = generateModifyTableSQL(tableName, dbTableDescription, structTableDescription)
	} else {
		rawSql = generateCreateTableSql(tableName, structTableDescription)
	}

	if rawSql != "" {
		if exec {
			_, err = db.Db.Exec(rawSql)
			if err != nil {
				return err
			}

			for _, v := range constraints {
				_, err = db.Db.Exec(v)
				if err != nil {
					return err
				}
			}
		} else {
			log.Print(rawSql)
			for _, v := range constraints {
				log.Print(v)
			}
		}
		if err != nil {
			return err
		} else if doLog {
			if success {
				log.Println("Updated table " + tableName)
			} else {
				log.Println("Created table " + tableName)
			}
		}
	}
	return nil
}

func generateModifyTableSQL(tableName string, dbTableDescription helpers.TableDescription, structTableDescription helpers.TableDescription) string {
	var columnsToAdd helpers.TableDescription
	var columnsToUpdate helpers.TableDescription
	var columnsToRemove helpers.TableDescription

	for i := range dbTableDescription {
		tableColumnDescription := dbTableDescription[i]
		structColumnDescription, succeeded := helpers.ArrFindFunc(
			structTableDescription,
			func(item helpers.ColumnDescription) bool {
				return item.Field == tableColumnDescription.Field
			},
		)
		if succeeded {
			if !tableColumnDescription.EqualTo(structColumnDescription) {
				columnsToUpdate = append(columnsToUpdate, structColumnDescription)
			}
		} else {
			columnsToRemove = append(columnsToRemove, tableColumnDescription)
		}
	}

	for i := range structTableDescription {
		structColumnDescription := structTableDescription[i]
		_, succeeded := helpers.ArrFindFunc(
			dbTableDescription,
			func(item helpers.ColumnDescription) bool {
				return item.Field == structColumnDescription.Field
			},
		)

		if !succeeded {
			columnsToAdd = append(columnsToAdd, structColumnDescription)
		}
	}

	if len(columnsToAdd) == 0 && len(columnsToUpdate) == 0 && len(columnsToRemove) == 0 {
		return ""
	}

	rawSql := "ALTER" + " TABLE " + tableName

	for i := range columnsToAdd {
		column := columnsToAdd[i]
		rawSql += " ADD COLUMN " + column.FormatSqlColumn() + ","
	}

	for i := range columnsToRemove {
		column := columnsToRemove[i]
		rawSql += " DROP COLUMN " + column.Field + ","
	}

	for i := range columnsToUpdate {
		column := columnsToUpdate[i]
		rawSql += " MODIFY COLUMN " + column.FormatSqlColumn() + ","
	}

	if rawSql[len(rawSql)-1:] == "," {
		rawSql = rawSql[:len(rawSql)-1]
	}

	return rawSql
}

func generateCreateTableSql(tableName string, structTableDescription helpers.TableDescription) string {
	rawSql := "CREATE TABLE " + tableName + "("

	var columns []string
	var constraints []string

	for i := range structTableDescription {
		column := structTableDescription[i]
		columns = append(columns, column.FormatSqlColumn())
		constraints = append(constraints, column.FormatSqlConstraints(tableName)...)
	}

	rawSql += strings.Join(columns, ", ") + strings.Join(constraints, ", ") + ")"

	return rawSql
}

func getTableDescription(db *sql.DB, tableName string) (helpers.TableDescription, bool, error) {
	var tableDescription helpers.TableDescription

	results, err := db.Query("DESCRIBE " + tableName)

	if err == nil {
		for results.Next() {
			var columnDescription helpers.ColumnDescription
			err = results.Scan(&columnDescription.Field, &columnDescription.Type, &columnDescription.Null, &columnDescription.Key, &columnDescription.Default, &columnDescription.Extra)

			if err != nil {
				return tableDescription, false, err
			}
			tableDescription = append(tableDescription, columnDescription)
		}

		return tableDescription, true, nil
	}

	return tableDescription, false, nil
}

func BuildTables(db *Database, entities ...interface{}) error {
	tableQueries, addConstraintsQueries, dropConstraintsQueries, err := GetUpdateTableQueriesForEntities(db, entities...)
	if err != nil {
		return err
	}

	for _, query := range dropConstraintsQueries {
		_, err = db.Db.Exec(query)
		if err != nil {
			return err
		}
	}

	for _, query := range tableQueries {
		_, err = db.Db.Exec(query)
		if err != nil {
			return err
		}
	}

	for _, query := range addConstraintsQueries {
		_, err = db.Db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
