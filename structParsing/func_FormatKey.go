package structParsing

import "strings"

// FormatKey takes the db_key struct tag string and determines if the column is a PRIMARY KEY column.
func FormatKey(key string) string {
	if strings.ToLower(key) == "primary" || strings.ToLower(key) == "pri" {
		return "PRIMARY KEY"
	}

	return ""
}
