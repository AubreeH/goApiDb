package helpers

import (
	"database/sql"
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
