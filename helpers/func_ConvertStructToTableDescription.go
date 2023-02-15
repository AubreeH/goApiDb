package helpers

import (
	"database/sql"
	"reflect"
)

func ConvertStructToTableDescription(entity any) TableDescription {
	structType := reflect.TypeOf(entity)
	structValue := reflect.ValueOf(entity)

	tableDescription := getTableDescription(structType, structValue)

	return tableDescription
}

func getTableDescription(refType reflect.Type, refValue reflect.Value) TableDescription {
	var description TableDescription

	numField := refType.NumField()
	for i := 0; i < numField; i++ {
		field := refType.Field(i)
		value := refValue.Field(i)

		parseStructTag := field.Tag.Get("parse_struct")
		parseStruct := parseStructTag == "" || ParseBool(parseStructTag)

		if value.Kind() == reflect.Struct && parseStruct {
			description = append(description, getTableDescription(value.Type(), value)...)
		} else {
			if !ParseBool(field.Tag.Get("sql_ignore")) {
				description = append(description, getColumnDescription(field))
			}
		}
	}

	return description
}

func getColumnDescription(field reflect.StructField) ColumnDescription {
	name := field.Tag.Get("sql_name")
	dataType := field.Tag.Get("sql_type")
	nullable := field.Tag.Get("sql_nullable")
	key := field.Tag.Get("sql_key")
	extra := field.Tag.Get("sql_extras")
	defaultValue := field.Tag.Get("sql_default")
	constraint := field.Tag.Get("sql_constraint")

	return ColumnDescription{
		Field: name,
		Type:  dataType,
		Null:  nullable,
		Key:   key,
		Default: sql.NullString{
			String: defaultValue,
		},
		Extra:           extra,
		StructFieldName: field.Name,
		Constraint:      constraint,
	}
}
