package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AubreeH/goApiDb/entities"
	"github.com/AubreeH/goApiDb/helpers"
	"reflect"
	"strings"
)

const (
	tagSqlName                         = "sql_name"
	tagSqlType                         = "sql_type"
	tagSqlKey                          = "sql_key"
	tagSqlExtras                       = "sql_extras"
	tagSqlNullable                     = "sql_nullable"
	tagSqlDefault                      = "sql_default"
	tagSqlDisallowExternalModification = "sql_disallow_external_modification"
)

func GetTableSqlDescriptionFromEntity[TEntity any]() (TablDesc, error) {
	var model TEntity
	tableDescription := TablDesc{}

	refValue := reflect.ValueOf(model)
	refType := reflect.TypeOf(model)

	if refType.Kind() != reflect.Struct {
		return TablDesc{}, errors.New("provided type is not a struct")
	}

	tableInfo, err := entities.GetTableInfo(model)
	if err != nil {
		return TablDesc{}, err
	}

	tableDescription.Name = tableInfo.Name

	for i := 0; i < refValue.NumField(); i++ {
		field := refValue.Type().Field(i)
		if field.Type != reflect.TypeOf(entities.EntityBase{}) {
			colDesc := parseColumn(field)
			tableDescription.Columns = append(tableDescription.Columns, colDesc)
			tableDescription.Constraints = append(tableDescription.Constraints, colDesc.GetConstraints(tableDescription.Name)...)
		}
	}

	return tableDescription, nil
}

func GetTableSqlDescriptionFromDb(db *Database, tableName string) (TablDesc, error) {
	if tableName == "" {
		return TablDesc{}, errors.New("empty table name provided")
	}

	result, err := db.Db.Query(fmt.Sprintf("DESCRIBE %s", tableName))
	if err != nil {
		return TablDesc{}, err
	}

	tableDescription := TablDesc{Name: tableName}
	for result.Next() {
		colDesc := ColDesc{}

		var sqlDefault sql.NullString
		err = result.Scan(&colDesc.Name, &colDesc.Type, &colDesc.Nullable, &colDesc.Key, &sqlDefault, &colDesc.Extras)
		if err != nil {
			return TablDesc{}, err
		}
		colDesc.Default = sqlDefault.String

		err = colDesc.getKeyFromDb(db, tableName)
		if err != nil {
			return TablDesc{}, err
		}

		tableDescription.Columns = append(tableDescription.Columns, colDesc)
		tableDescription.Constraints = append(tableDescription.Constraints, colDesc.GetConstraints(tableName)...)
	}

	return tableDescription, nil
}

func parseColumn(structField reflect.StructField) ColDesc {
	desc := ColDesc{}
	helpers.TagLookup(structField, tagSqlName, &desc.Name)
	helpers.TagLookup(structField, tagSqlType, &desc.Type)
	helpers.TagLookup(structField, tagSqlKey, &desc.Key)
	helpers.TagLookup(structField, tagSqlExtras, &desc.Extras)
	helpers.TagLookup(structField, tagSqlNullable, &desc.Nullable)
	helpers.TagLookup(structField, tagSqlDefault, &desc.Default)

	var output string
	helpers.TagLookup(structField, tagSqlDisallowExternalModification, &output)
	desc.DisallowExternalModification = helpers.ParseBool(output)

	return desc
}

func GetUpdateTableQueries[TEntity any](db *Database) (tableSql string, addConstraintsSql []string, dropConstraintsSql []string, err error) {
	entityDesc, err := GetTableSqlDescriptionFromEntity[TEntity]()
	if err != nil {
		return "", nil, nil, err
	}

	dbDesc, err := GetTableSqlDescriptionFromDb(db, entityDesc.Name)
	if err != nil {
		if strings.Index(err.Error(), "Error 1146 (42S02)") == -1 {
			return "", nil, nil, err
		}

		tableSql, addConstraintsSql = entityDesc.Format()

		return tableSql, addConstraintsSql, nil, nil
	}

	diff, err := GetDescriptionDifferences(entityDesc, dbDesc)
	if err != nil {
		return "", nil, nil, err
	}

	tableSql, addConstraintsSql, dropConstraintsSql = diff.Format()

	return tableSql, addConstraintsSql, dropConstraintsSql, nil
}

func GetDescriptionDifferences(entityDesc TablDesc, dbDesc TablDesc) (TablDescDiff, error) {
	diff := TablDescDiff{}

	columnsToAdd, columnsToModify, columnsToDrop, err := GetColumnDifferences(entityDesc.Name, entityDesc.Columns, dbDesc.Columns)
	if err != nil {
		return TablDescDiff{}, err
	}

	diff.ColumnsToAdd = columnsToAdd
	diff.ColumnsToModify = columnsToModify
	diff.ColumnsToDrop = columnsToDrop

	constraintsToAdd, constraintsToDrop := GetConstraintDifferences(entityDesc.Constraints, dbDesc.Constraints)

	diff.ConstraintsToAdd = constraintsToAdd
	diff.ConstraintsToDrop = constraintsToDrop

	return diff, nil
}

func GetColumnDifferences(tableName string, entityColumns []ColDesc, dbColumns []ColDesc) (add []ColDesc, modify []ColDesc, drop []ColDesc, err error) {
	add = []ColDesc{}
	modify = []ColDesc{}
	drop = []ColDesc{}

	for _, entityCol := range entityColumns {
		if entityCol.Name == "" {
			return nil, nil, nil, errors.New("column with empty name provided with entity description")
		}

		dbCol, ok := helpers.ArrFindFunc(dbColumns, func(dbCol ColDesc) bool {
			return entityCol.Name == dbCol.Name
		})
		if !ok {
			add = append(add, entityCol)
		} else {
			entityColSql := entityCol.Format(tableName)
			dbColSql := dbCol.Format(tableName)

			if entityColSql != dbColSql {
				modify = append(modify, entityCol)
			}
		}
	}

	for _, dbCol := range dbColumns {
		if dbCol.Name == "" {
			return nil, nil, nil, errors.New("column with empty name provided with db description")
		}

		_, ok := helpers.ArrFindFunc(entityColumns, func(entityCol ColDesc) bool {
			return entityCol.Name == dbCol.Name
		})

		if !ok {
			drop = append(drop, dbCol)
		}
	}

	return add, modify, drop, nil
}

func GetConstraintDifferences(entityConstraints []Constraint, dbConstraints []Constraint) (add []Constraint, drop []Constraint) {
	add = []Constraint{}
	drop = []Constraint{}

	for _, entityConstraint := range entityConstraints {
		_, ok := helpers.ArrFindFunc(dbConstraints, func(dbConstraint Constraint) bool {
			return entityConstraint.TableName == dbConstraint.TableName &&
				entityConstraint.ColumnName == dbConstraint.ColumnName &&
				entityConstraint.ReferencedTableName == dbConstraint.ReferencedTableName &&
				entityConstraint.ReferencedColumnName == dbConstraint.ReferencedColumnName
		})

		if !ok {
			add = append(add, entityConstraint)
		}
	}

	for _, dbConstraint := range dbConstraints {
		_, ok := helpers.ArrFindFunc(entityConstraints, func(entityConstraint Constraint) bool {
			return entityConstraint.TableName == dbConstraint.TableName &&
				entityConstraint.ColumnName == dbConstraint.ColumnName &&
				entityConstraint.ReferencedTableName == dbConstraint.ReferencedTableName &&
				entityConstraint.ReferencedColumnName == dbConstraint.ReferencedColumnName
		})

		if !ok {
			drop = append(drop, dbConstraint)
		}
	}

	return add, drop
}
