package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/AubreeH/goApiDb/structParsing"
)

func GetTableSqlDescriptionFromDb(db *Database, tableName string) (structParsing.TablDesc, error) {
	if tableName == "" {
		return structParsing.TablDesc{}, errors.New("empty table name provided")
	}

	result, err := db.Db.Query(fmt.Sprintf("DESCRIBE %s", tableName))
	if err != nil {
		return structParsing.TablDesc{}, err
	}

	defer result.Close()

	tableDescription := structParsing.TablDesc{Name: tableName}
	for result.Next() {
		colDesc := structParsing.ColDesc{}

		var sqlDefault sql.NullString
		err = result.Scan(&colDesc.Name, &colDesc.Type, &colDesc.Nullable, &colDesc.Key, &sqlDefault, &colDesc.Extras)
		if err != nil {
			return structParsing.TablDesc{}, err
		}
		colDesc.Type = structParsing.FormatType(colDesc.Type)
		colDesc.Default = structParsing.FormatDefault(sqlDefault.String)
		colDesc.Nullable = structParsing.FormatNullable(colDesc.Nullable)
		colDesc.Extras = structParsing.FormatExtras(colDesc.Extras)

		err = getKeyFromDb(db, tableName, &colDesc)
		if err != nil {
			return structParsing.TablDesc{}, err
		}

		tableDescription.Columns = append(tableDescription.Columns, colDesc)
		tableDescription.Constraints = append(tableDescription.Constraints, colDesc.GetConstraints(tableName)...)
	}

	return tableDescription, nil
}

func GetUpdateTableQueriesForEntity[TEntity any](db *Database, entity TEntity) (tableSql string, addConstraintsSql []string, dropConstraintsSql []string, err error) {
	entityDesc, err := structParsing.GetTableSqlDescriptionFromEntity(entity)
	if err != nil {
		return "", nil, nil, err
	}

	dbDesc, err := GetTableSqlDescriptionFromDb(db, entityDesc.Name)
	if err != nil {
		if !strings.Contains(err.Error(), "Error 1146 (42S02)") {
			return "", nil, nil, err
		}

		tableSql, addConstraintsSql = entityDesc.Format()

		return tableSql, addConstraintsSql, nil, nil
	}

	diff, err := getDescriptionDifferences(entityDesc, dbDesc)
	if err != nil {
		return "", nil, nil, err
	}

	tableSql, addConstraintsSql, dropConstraintsSql = diff.Format()

	return tableSql, addConstraintsSql, dropConstraintsSql, nil
}

func GetUpdateTableQueriesForEntities(db *Database, entities ...interface{}) (tableQueries []string, addConstraintsQueries []string, dropConstraintsQueries []string, err error) {
	tableQueries = []string{}
	addConstraintsQueries = []string{}
	dropConstraintsQueries = []string{}

	for _, e := range entities {
		tableSql, addConstraintsSql, dropConstraintsSql, err := GetUpdateTableQueriesForEntity(db, e)
		if err != nil {
			return nil, nil, nil, err
		}
		if tableSql != "" {
			tableQueries = append(tableQueries, tableSql)
		}
		addConstraintsQueries = append(addConstraintsQueries, addConstraintsSql...)
		dropConstraintsQueries = append(dropConstraintsQueries, dropConstraintsSql...)
	}

	return tableQueries, addConstraintsQueries, dropConstraintsQueries, nil
}

// getKeyFromDb retrieves the key for the column currently defined in the db.
// TODO: Add unique constraint
func getKeyFromDb(db *Database, tableName string, col *structParsing.ColDesc) error {
	if col.Name == "" {
		return errors.New("column name is empty")
	}

	result, err := db.Db.Query(`SELECT CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = ?`, db.dbName, tableName, col.Name)
	if err != nil {
		return err
	}

	defer result.Close()

	if result.Next() {
		var constraintName string
		var columnName string
		var referencedTableName sql.NullString
		var referencedColumnName sql.NullString
		err = result.Scan(&constraintName, &columnName, &referencedTableName, &referencedColumnName)
		if err != nil {
			return err
		}

		if strings.ToLower(constraintName) == "primary" {
			col.Key = "primary"
		} else if referencedTableName.Valid && referencedColumnName.Valid {
			col.Key = fmt.Sprintf("foreign,%s,%s,%s", referencedTableName.String, referencedColumnName.String, constraintName)
		}
	}

	return nil
}
