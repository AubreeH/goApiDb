package structParsing

import "strings"

func (col *ColDesc) parseType() string {
	switch strings.ToUpper(col.Type) {
	case "BOOL", "BOOLEAN", "TINYINT(1)":
		return "BOOLEAN"
	case "TIMESTAMP", "DATETIME":
		return "DATETIME"
	case "INT", "INT(11)":
		return "INTEGER"
	default:
		return col.Type
	}
}
