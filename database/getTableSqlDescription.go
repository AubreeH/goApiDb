package database

import (
	"errors"
	"fmt"
	"github.com/AubreeH/goApiDb/helpers"
	"log"
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

type TablDesc struct {
	Name    string
	Columns []ColDesc
}

type ColDesc struct {
	Name                         string
	Type                         string
	Key                          string
	Extras                       string
	Nullable                     string
	Default                      string
	DisallowExternalModification bool
}

func GetTableSqlDescription[TModel any]() error {
	var model TModel

	refValue := reflect.ValueOf(model)
	refType := reflect.TypeOf(model)

	if refType.Kind() != reflect.Struct {
		return errors.New("provided type is not a struct")
	}

	var columns []string
	var constraints []string

	for i := 0; i < refValue.NumField(); i++ {
		colString, colConstraints := parseColumn(helpers.GetTableName(model), refValue.Type().Field(i))
		columns = append(columns, colString)
		constraints = append(constraints, colConstraints...)
	}

	log.Print(columns, constraints)

	return errors.New("not yet implemented")
}

func parseColumn(tableName string, structField reflect.StructField) (colString string, constraints []string) {
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

	return desc.Format(tableName)
}

func (col ColDesc) Format(tableName string) (string, []string) {
	var s []string
	var constraints []string

	key, keyConstraint := formatKey(tableName, col.Name, col.Key)
	extras := formatExtras(col.Extras)
	nullable := formatNullable(col.Nullable)
	def := formatDefault(col.Default)
	t := formatType(col.Type)

	helpers.ArrAdd(&constraints, keyConstraint)
	helpers.ArrAdd(&s, col.Name, t, key, nullable, def, extras)

	return strings.Join(s, " "), constraints
}

func formatKey(tableName string, columnName string, key string) (out string, constraint string) {
	if key == "" {
		return "", ""
	}

	s := strings.Split(key, ",")
	if strings.ToLower(s[0]) == "primary" {
		return "PRIMARY KEY", ""
	} else if strings.ToLower(s[0]) == "foreign" {
		if len(s) != 3 {
			return "", ""
		}

		fkName := fmt.Sprintf("FK_%s_%s_%s_%s", tableName, columnName, s[1], s[2])
		c := fmt.Sprintf("ALTER TABLE %s ADD FOREIGN KEY %s REFERENCES %s(%s)", tableName, fkName, s[1], s[2])

		return "", c
	}

	return "", ""
}

func formatExtras(extras string) string {
	return extras
}

func formatNullable(nullable string) string {
	if helpers.ParseBool(nullable) {
		return ""
	}

	return "NOT NULL"
}

func formatDefault(def string) string {
	if def == "" {
		return ""
	}

	return "DEFAULT " + def
}

func formatType(t string) string {
	return t
}
