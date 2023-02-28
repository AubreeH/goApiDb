package database

import (
	"github.com/AubreeH/goApiDb/helpers"
	"strings"
)

func formatKey(key string) string {
	if strings.ToLower(key) == "primary" || strings.ToLower(key) == "pri" {
		return "PRIMARY KEY"
	}

	return ""
}

func formatExtras(extras string) string {
	return strings.ToUpper(extras)
}

func formatNullable(nullable string) string {
	if helpers.ParseBool(nullable) {
		return ""
	}

	return "NOT NULL"
}

func formatDefault(def string) string {
	if def == "" {
		return ""
	}

	return "DEFAULT " + def
}

func formatType(t string) string {
	return strings.ToUpper(t)
}
