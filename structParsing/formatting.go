package structParsing

import (
	"github.com/AubreeH/goApiDb/helpers"
	"strings"
)

// FormatName takes the db_name struct tag string and the name of the field and returns either the db_name tag if it is set or the field name in snake case.
func FormatName(tag string, fieldName string) string {
	if tag != "" {
		return tag
	}

	return helpers.PascalToSnakeCase(fieldName)
}

// FormatKey takes the db_key struct tag string and determines if the column is a PRIMARY KEY column.
func FormatKey(key string) string {
	if strings.ToLower(key) == "primary" || strings.ToLower(key) == "pri" {
		return "PRIMARY KEY"
	}

	return ""
}

// FormatExtras takes the db_extras struct tag string and formats it for use in database manipulation.
func FormatExtras(extras string) string {
	return strings.ToUpper(extras)
}

// FormatNullable takes the db_null struct tag string and determines if the column is nullable or not.
func FormatNullable(nullable string) string {
	if FormatBoolean(nullable) != 1 {
		return "NOT NULL"
	}

	return ""
}

// FormatDefault takes the db_default struct tag string and determines the default values of the field.
func FormatDefault(def string) string {
	if def == "" {
		return ""
	}

	return "DEFAULT " + def
}

// FormatType takes the db_type struct tag string and determines the SQL type for the field.
func FormatType(t string) string {
	return strings.ToUpper(t)
}

// FormatBoolean takes any boolean tag values and returns 0 for false, 1 for true and -1 if no value was provided/could be interpreted
func FormatBoolean(b string) int {
	switch strings.ToLower(b) {
	case "true":
		return 1
	case "yes":
		return 1
	case "y":
		return 1
	case "false":
		return 0
	case "no":
		return 0
	case "n":
		return 0
	default:
		return -1
	}
}
