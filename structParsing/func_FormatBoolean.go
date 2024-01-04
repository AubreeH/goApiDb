package structParsing

import "strings"

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
