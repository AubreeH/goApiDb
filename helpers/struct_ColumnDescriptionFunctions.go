package helpers

import (
	"fmt"
	"regexp"
)

type SqlKey interface {
	string
}

func (description ColumnDescription) EqualTo(j ColumnDescription) bool {
	if description.Field != j.Field {
		return false
	}
	if ParseColumnDescriptionType(description.Type) != ParseColumnDescriptionType(j.Type) {
		return false
	}
	if ParseColumnDescriptionNullable(description.Null) != ParseColumnDescriptionNullable(j.Null) {
		return false
	}
	if ParseColumnDescriptionKey(description.Key) != ParseColumnDescriptionKey(j.Key) {
		return false
	}
	if description.Default != j.Default {
		return false
	}
	if ParseColumnDescriptionExtra(description.Extra) != ParseColumnDescriptionExtra(j.Extra) {
		return false
	}
	return true
}

func (description ColumnDescription) FormatSqlColumn() string {
	var sqlString string
	sqlString += description.Field + " "
	sqlString += ParseColumnDescriptionType(description.Type) + " "
	sqlString += ParseColumnDescriptionKey(description.Key) + " "
	sqlString += ParseColumnDescriptionNullable(description.Null) + " "
	if description.Default.Valid && description.Default.String != "" {
		sqlString += "DEFAULT '" + description.Default.String + "' "
	}
	sqlString += description.Extra
	if description.Constraint != "" {
		sqlString += "CHECK (" + description.Constraint + ")"
	}

	return sqlString
}

func (description ColumnDescription) FormatSqlConstraints(tableName string) []string {
	var output []string

	re, err := regexp.Compile("foreign,(?P<table>.*?),(?P<column>.*)")
	if err != nil {
		panic(err)
	}

	submatch := re.FindStringSubmatch(description.Key)
	if len(submatch) > 1 {
		constraint := fmt.Sprintf("ALTER TABLE %s ADD FOREIGN KEY (%s) REFERENCES %s(%s)", tableName, description.Field, submatch[1], submatch[2])
		output = append(output, constraint)
	}

	return output
}
