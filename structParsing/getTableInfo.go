package structParsing

import (
	"errors"
	"reflect"

	"github.com/AubreeH/goApiDb/entities"
	"github.com/AubreeH/goApiDb/helpers"
)

type TableInfo struct {
	Name        string
	SoftDeletes string
	IsValid     bool
	PrimaryKey  string
}

var entityBaseType = reflect.TypeOf(entities.EntityBase{})

func GetTableInfo(entity interface{}) (TableInfo, error) {
	entityVal := reflect.ValueOf(entity)

	if entityVal.Kind() == reflect.Ptr || entityVal.Kind() == reflect.Interface {
		entityVal = entityVal.Elem()

		if entityVal.Kind() == reflect.Interface {
			entityVal = entityVal.Elem()
		}
	}

	entityType := entityVal.Type()

	if entityType.Kind() != reflect.Struct {
		return TableInfo{}, errors.New("provided entity is not a struct")
	}

	tableInfo := TableInfo{}

	if tableInfo.Name == "" {
		tableInfo.Name = helpers.GetTableName(entity)
	}

	getInfo(&tableInfo, entityVal, entityType)

	if !tableInfo.IsValid {
		return TableInfo{}, errors.New("no entity base in struct")
	}

	tableInfo.PrimaryKey = getPrimaryKey(entityVal, entityType)

	return tableInfo, nil
}

func getInfo(tableInfo *TableInfo, entityValue reflect.Value, baseType reflect.Type) {
	entityType := entityValue.Type()
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		fieldValue := entityValue.Field(i)

		if field.Type == entityBaseType {
			tableInfo.Name = FormatName(field.Tag.Get("table_name"), baseType.Name())
			tableInfo.IsValid = true
		} else if fieldValue.Kind() == reflect.Struct && FormatParseStruct(field) {
			getInfo(tableInfo, fieldValue, baseType)
		} else {
			if FormatSoftDeletes(field) {
				tableInfo.SoftDeletes = FormatSqlName(field)
			}
		}
	}
}

func getPrimaryKey(entityValue reflect.Value, baseType reflect.Type) string {
	entityType := entityValue.Type()
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		fieldValue := entityValue.Field(i)

		if FormatKey(field.Tag.Get(SqlKey)) == "PRIMARY KEY" {
			return FormatSqlName(field)
		} else if fieldValue.Kind() == reflect.Struct && FormatParseStruct(field) {
			pk := getPrimaryKey(fieldValue, baseType)
			if pk != "" {
				return pk
			}
		}
	}
	return ""
}
