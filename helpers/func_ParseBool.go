package helpers

import "strings"

func ParseBool(value string) bool {
	lowerCaseValue := strings.ToLower(value)

	switch lowerCaseValue {
	case "true":
		return true
	case "yes":
		return true
	case "nullable":
		return true
	case "false":
		return false
	case "no":
		return false
	case "not null":
		return false
	case "":
		return false
	}

	return false
}
