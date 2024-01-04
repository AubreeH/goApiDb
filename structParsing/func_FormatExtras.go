package structParsing

import "strings"

// FormatExtras takes the db_extras struct tag string and formats it for use in database manipulation.
func FormatExtras(extras string) string {
	return strings.ToUpper(extras)
}
