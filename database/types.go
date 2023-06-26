package database

import (
	"fmt"
	"strings"

	"github.com/AubreeH/goApiDb/structParsing"
)

type TablDescDiff struct {
	Name              string
	ColumnsToAdd      []structParsing.ColDesc
	ColumnsToModify   []structParsing.ColDesc
	ColumnsToDrop     []structParsing.ColDesc
	ConstraintsToAdd  []structParsing.Constraint
	ConstraintsToDrop []structParsing.Constraint
}

func (diff TablDescDiff) Format() (tableQuery string, addConstraintQueries []string, dropConstraintQueries []string) {
	addConstraintQueries = []string{}
	dropConstraintQueries = []string{}
	columns := ""

	for _, col := range diff.ColumnsToAdd {
		colSql := col.Format(diff.Name)
		if columns == "" {
			columns += fmt.Sprintf("ADD %s", colSql)
		} else {
			columns += fmt.Sprintf(", ADD %s", colSql)
		}
	}

	for _, col := range diff.ColumnsToModify {
		colSql := col.Format(diff.Name)
		if columns == "" {
			columns += fmt.Sprintf("MODIFY %s", colSql)
		} else {
			columns += fmt.Sprintf(", MODIFY %s", colSql)
		}
	}

	for _, col := range diff.ColumnsToDrop {
		colSql := col.Format(diff.Name, "drop")
		if columns == "" {
			columns += fmt.Sprintf("DROP %s", colSql)
		} else {
			columns += fmt.Sprintf(", DROP %s", colSql)
		}
	}

	for _, constraint := range diff.ConstraintsToDrop {
		dropConstraintQueries = append(dropConstraintQueries, constraint.Format("drop"))
	}

	for _, constraint := range diff.ConstraintsToAdd {
		addConstraintQueries = append(addConstraintQueries, constraint.Format("add"))
	}

	columns = strings.TrimSpace(columns)
	if columns == "" {
		return "", addConstraintQueries, dropConstraintQueries
	}

	tableQuery = fmt.Sprintf("ALTER TABLE `%s` %s;", diff.Name, columns)
	return tableQuery, addConstraintQueries, dropConstraintQueries
}
