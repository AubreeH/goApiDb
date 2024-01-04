package structParsing

import "github.com/AubreeH/goApiDb/helpers"

// FormatName takes the db_name struct tag string and the name of the field and returns either the db_name tag if it is set or the field name in snake case.
func FormatName(tag string, fieldName string) string {
	if tag != "" {
		return tag
	}

	return helpers.PascalToSnakeCase(fieldName)
}
