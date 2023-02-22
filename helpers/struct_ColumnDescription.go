package helpers

import (
	"database/sql"
	"log"
	"regexp"
	"strings"
)

type TableDescription []ColumnDescription

type ColumnDescription struct {
	Field      string
	Type       string
	Null       string
	Key        string
	Default    sql.NullString
	Extra      string
	Constraint string

	StructFieldName string
}

func (description ColumnDescription) EqualTo(j ColumnDescription) bool {
	if description.Field != j.Field {
		return false
	}
	if ParseColumnDescriptionType(description.Type) != ParseColumnDescriptionType(j.Type) {
		return false
	}
	if ParseColumnDescriptionNullable(description.Null) != ParseColumnDescriptionNullable(j.Null) {
		return false
	}
	if ParseColumnDescriptionKey(description.Key) != ParseColumnDescriptionKey(j.Key) {
		return false
	}
	if description.Default != j.Default {
		return false
	}
	if ParseColumnDescriptionExtra(description.Extra) != ParseColumnDescriptionExtra(j.Extra) {
		return false
	}
	return true
}

func ParseColumnDescriptionType(value string) string {
	lowerCaseValue := strings.ToLower(value)

	switch lowerCaseValue {
	case "int":
		return "int(11)"
	case "string":
		return "varchar(128)"
	}

	return lowerCaseValue
}

func ParseColumnDescriptionNullable(value string) string {
	lowerCaseValue := strings.ToLower(value)

	switch lowerCaseValue {
	case "true":
		return "NULL"
	case "yes":
		return "NULL"
	case "nullable":
		return "NULL"
	case "false":
		return "NOT NULL"
	case "no":
		return "NOT NULL"
	case "not null":
		return "NOT NULL"
	}

	return lowerCaseValue
}

func ParseColumnDescriptionKey(value string) string {
	lowerCaseValue := strings.ToLower(value)

	switch lowerCaseValue {
	case "pri":
		return "PRIMARY KEY"
	case "primary":
		return "PRIMARY KEY"
	}

	return ""
}

func ParseColumnDescriptionExtra(value string) string {
	return strings.ToLower(value)
}

func (description ColumnDescription) FormatSqlColumn() string {
	var sqlString string
	sqlString += description.Field + " "
	sqlString += ParseColumnDescriptionType(description.Type) + " "
	sqlString += ParseColumnDescriptionKey(description.Key) + " "
	sqlString += ParseColumnDescriptionNullable(description.Null) + " "
	if description.Default.Valid && description.Default.String != "" {
		sqlString += "DEFAULT '" + description.Default.String + "' "
	}
	sqlString += description.Extra
	if description.Constraint != "" {
		sqlString += "CHECK (" + description.Constraint + ")"
	}

	return sqlString
}

func (description ColumnDescription) FormatSqlConstraints(tableName string) []string {
	var output []string

	re, err := regexp.Compile("foreign,(?P<table>.*?),(?P<column>.*)")
	if err != nil {
		panic(err)
	}

	submatch := re.FindStringSubmatch(description.Key)
	if len(submatch) > 1 {
		constraint := "ALTER TABLE " + tableName + " ADD FOREIGN KEY (" + description.Field + ") REFERENCES " + submatch[1] + "(" + submatch[2] + ")"
		log.Print(constraint)
	}

	return output
}
