package structParsing

import "reflect"

const (
	DbName                         Tag = "db_name"
	DbType                         Tag = "db_type"
	DbKey                          Tag = "db_key"
	DbExtras                       Tag = "db_extras"
	DbNullable                     Tag = "db_nullable"
	DbDefault                      Tag = "db_default"
	DbDisallowExternalModification Tag = "db_disallow_external_modification"
	DbIgnore                       Tag = "db_ignore"
	DbParseStruct                  Tag = "parse_struct"
	DbSoftDeletes                  Tag = "soft_deletes"
)

func FormatSqlName(field reflect.StructField) string {
	return FormatName(DbName.Get(field), field.Name)
}

func FormatSqlType(field reflect.StructField) string {
	return FormatType(DbType.Get(field))
}

func FormatSqlKey(field reflect.StructField) string {
	return FormatKey(DbKey.Get(field))
}

func FormatSqlExtras(field reflect.StructField) string {
	return FormatExtras(DbExtras.Get(field))
}

func FormatSqlNullable(field reflect.StructField) string {
	return FormatNullable(DbNullable.Get(field))
}

func FormatSqlDefault(field reflect.StructField) string {
	return FormatDefault(DbDefault.Get(field))
}

func FormatSqlDisallowExternalModification(field reflect.StructField) bool {
	disallowExternalModVal := DbDisallowExternalModification.Get(field)
	keyVal := FormatKey(DbKey.Get(field))

	if disallowExternalModVal == "" && keyVal == "PRIMARY KEY" {
		return true
	}

	return FormatBoolean(disallowExternalModVal) == 1
}

func FormatSqlIgnore(field reflect.StructField) bool {
	return FormatBoolean(DbIgnore.Get(field)) == 1
}

func FormatParseStruct(field reflect.StructField) bool {
	return FormatBoolean(DbParseStruct.Get(field)) != 0
}

func FormatSoftDeletes(field reflect.StructField) bool {
	return FormatBoolean(DbSoftDeletes.Get(field)) == 1
}
