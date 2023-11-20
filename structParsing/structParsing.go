package structParsing

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/AubreeH/goApiDb/helpers"
)

func GetTableSqlDescriptionFromEntity[TEntity interface{}](entity TEntity) (tablDesc TablDesc, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("an error occurred whilst parsing table description from provided entity: %v", err)
		}
	}()

	tableDescription := TablDesc{}

	refValue := helpers.GetRootValue(reflect.ValueOf(entity))
	refType := refValue.Type()

	if !refValue.IsValid() {
		return TablDesc{}, errors.New("this value is invalid")
	}

	if refType.Kind() != reflect.Struct {
		return TablDesc{}, errors.New("provided type is not a struct")
	}

	tableInfo, err := GetTableInfo(entity)
	if err != nil {
		return TablDesc{}, err
	}

	tableDescription.Name = tableInfo.Name

	parseColumns(&tableDescription, refValue)

	return tableDescription, nil
}

func parseColumns(tableDescription *TablDesc, refValue reflect.Value) {
	for i := 0; i < refValue.NumField(); i++ {
		field := refValue.Type().Field(i)
		if field.Type != reflect.TypeOf(entityBaseType) && !FormatSqlIgnore(field) {
			if field.Type.Kind() == reflect.Struct && FormatParseStruct(field) {
				parseColumns(tableDescription, refValue.Field(i))
			} else {
				colDesc := parseColumn(field, refValue.Field(i))
				tableDescription.Columns = append(tableDescription.Columns, colDesc)
				tableDescription.Constraints = append(tableDescription.Constraints, colDesc.GetConstraints(tableDescription.Name)...)
			}
		}
	}
}

func parseColumn(structField reflect.StructField, fieldValue reflect.Value) ColDesc {
	desc := ColDesc{}
	desc.Type = FormatSqlType(structField)
	desc.Key = GetTag(structField, SqlKey)
	desc.Extras = FormatSqlExtras(structField)
	desc.Nullable = FormatSqlNullable(structField)
	desc.Default = FormatSqlDefault(structField)

	if GetTag(structField, SqlDisallowExternalModification) == "" && FormatKey(desc.Key) == "PRIMARY KEY" {
		desc.DisallowExternalModification = true
	} else {
		desc.DisallowExternalModification = FormatSqlDisallowExternalModification(structField)
	}

	var name string
	helpers.TagLookup(structField, SqlName, &name)
	desc.Name = FormatName(name, structField.Name)

	if fieldValue.CanAddr() {
		desc.Pointer = fieldValue.Addr().Interface()
	}

	desc.Value = fieldValue.Interface()

	return desc
}
