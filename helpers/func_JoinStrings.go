package helpers

import "strings"

func JoinStrings(seperator string, s ...string) string {
	return strings.Join(s, seperator)
}
