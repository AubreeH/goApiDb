package structParsing

import (
	"errors"
	"github.com/AubreeH/goApiDb/entities"
	"github.com/AubreeH/goApiDb/helpers"
	"reflect"
)

type TableInfo struct {
	Name        string
	SoftDeletes string
	IsValid     bool
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