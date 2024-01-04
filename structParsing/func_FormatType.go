package structParsing

import "strings"

// FormatType takes the db_type struct tag string and determines the SQL type for the field.
func FormatType(t string) string {
	return strings.ToUpper(t)
}
