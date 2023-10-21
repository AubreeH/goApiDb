package structParsing

import (
	"strings"
)

func (col *ColDesc) parseDefaultValue() string {
	switch strings.ToUpper(col.Type) {
	case "BOOL":
		return col.parseDefaultValueBool()
	case "DATETIME":
		return col.parseDefaultValueDateTime()
	default:
		return col.Default
	}
}

func (col *ColDesc) parseDefaultValueBool() string {
	switch strings.ToUpper(col.Default) {
	case "1":
		return "TRUE"
	case "0":
		return "FALSE"
	default:
		return col.Default
	}
}

func (col *ColDesc) parseDefaultValueDateTime() string {
	switch strings.ToUpper(col.Default) {
	case "DEFAULT CURRENT_TIMESTAMP", "DEFAULT CURRENT_TIMESTAMP()":
		return "DEFAULT CURRENT_TIMESTAMP"
	default:
		return col.Default
	}
}
