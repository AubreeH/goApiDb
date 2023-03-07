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

func GetTableSqlDescriptionFromEntity[TEntity interface{}](entity TEntity) (TablDesc, error) {
	tableDescription := TablDesc{}

	refValue := reflect.ValueOf(entity)
	refType := refValue.Type()

	if !refValue.IsValid() {
		return TablDesc{}, errors.New("this value is invalid")
	}

	if refType.Kind() == reflect.Interface {
		refValue = refValue.Elem()
		refType = refValue.Type()
	}

	if refType.Kind() != reflect.Struct {
		return TablDesc{}, errors.New("provided type is not a struct")
	}

	tableInfo, err := entities.GetTableInfo(entity)
	if err != nil {
		return TablDesc{}, err
	}

	tableDescription.Name = tableInfo.Name

	for i := 0; i < refValue.NumField(); i++ {
		field := refValue.Type().Field(i)
		if field.Type != reflect.TypeOf(entities.EntityBase{}) && !helpers.ParseBool(field.Tag.Get("sql_ignore")) {
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

func GetUpdateTableQueriesForEntity[TEntity any](db *Database, entity TEntity) (tableSql string, addConstraintsSql []string, dropConstraintsSql []string, err error) {
	entityDesc, err := GetTableSqlDescriptionFromEntity(entity)
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
		tableQueries = append(tableQueries, tableSql)
		addConstraintsQueries = append(addConstraintsQueries, addConstraintsSql...)
		dropConstraintsQueries = append(dropConstraintsQueries, dropConstraintsSql...)
	}

	return tableQueries, addConstraintsQueries, dropConstraintsQueries, nil
}
