package structParsing

import (
	"reflect"
)

const (
	SqlName                         = "db_name"
	SqlType                         = "db_type"
	SqlKey                          = "db_key"
	SqlExtras                       = "db_extras"
	SqlNullable                     = "db_nullable"
	SqlDefault                      = "db_default"
	SqlDisallowExternalModification = "db_disallow_external_modification"
	SqlIgnore                       = "db_ignore"
	ParseStruct                     = "parse_struct"
	SoftDeletes                     = "soft_deletes"
)

func GetTag(structField reflect.StructField, tag string) string {
	return structField.Tag.Get(tag)
}

func FormatSqlName(structField reflect.StructField) string {
	return FormatName(GetTag(structField, SqlName), structField.Name)
}

func FormatSqlType(structField reflect.StructField) string {
	return FormatType(GetTag(structField, SqlType))
}

func FormatSqlKey(structField reflect.StructField) string {
	return FormatKey(GetTag(structField, SqlKey))
}

func FormatSqlExtras(structField reflect.StructField) string {
	return FormatExtras(GetTag(structField, SqlExtras))
}

func FormatSqlNullable(structField reflect.StructField) string {
	return FormatNullable(GetTag(structField, SqlNullable))
}

func FormatSqlDefault(structField reflect.StructField) string {
	return FormatDefault(GetTag(structField, SqlDefault))
}

func FormatSqlDisallowExternalModification(structField reflect.StructField) bool {
	return FormatBoolean(GetTag(structField, SqlDisallowExternalModification)) == 1
}

func FormatSqlIgnore(structField reflect.StructField) bool {
	return FormatBoolean(GetTag(structField, SqlIgnore)) == 1
}

func FormatParseStruct(structField reflect.StructField) bool {
	return FormatBoolean(GetTag(structField, ParseStruct)) != 0
}

func FormatSoftDeletes(structField reflect.StructField) bool {
	return FormatBoolean(GetTag(structField, SoftDeletes)) == 1
}
