package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AubreeH/goApiDb/helpers"
	"log"
	"strings"
)

type TablDesc struct {
	Name        string
	Columns     []ColDesc
	Constraints []Constraint
}

type TablDescDiff struct {
	Name              string
	ColumnsToAdd      []ColDesc
	ColumnsToModify   []ColDesc
	ColumnsToDrop     []ColDesc
	ConstraintsToAdd  []Constraint
	ConstraintsToDrop []Constraint
}

type ColDesc struct {
	Name                         string
	Type                         string
	Key                          string
	Extras                       string
	Nullable                     string
	Default                      string
	DisallowExternalModification bool
}

type Constraint struct {
	ConstraintName       string
	TableName            string
	ColumnName           string
	ReferencedTableName  string
	ReferencedColumnName string
}

func (tabl TablDesc) Format() (string, []string) {

	var columns []string
	var constraints []string

	for _, col := range tabl.Columns {
		colString := col.Format(tabl.Name)
		columns = append(columns, colString)

		for _, v := range col.GetConstraints(tabl.Name) {
			constraints = append(constraints, v.Format("add"))
		}
	}

	log.Print(constraints)

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tabl.Name, strings.Join(columns, ", ")), constraints
}

func (col *ColDesc) Format(tableName string) string {
	var s []string

	key := formatKey(col.Key)
	extras := formatExtras(col.Extras)
	nullable := formatNullable(col.Nullable)
	def := formatDefault(col.Default)
	t := formatType(col.Type)

	helpers.ArrAdd(&s, col.Name, t, key, nullable, def, extras)

	return strings.Join(s, " ")
}

func (col *ColDesc) GetConstraints(tableName string) []Constraint {
	var constraints []Constraint

	s := strings.Split(col.Key, ",")
	if (len(s) == 3 || len(s) == 4) && strings.ToLower(s[0]) == "foreign" {
		var fkName string
		if len(s) == 4 {
			fkName = s[3]
		} else {
			fkName = fmt.Sprintf("FK_%s_%s_%s_%s", tableName, col.Name, s[1], s[2])
		}

		fk := Constraint{ConstraintName: fkName, TableName: tableName, ColumnName: col.Name, ReferencedTableName: s[1], ReferencedColumnName: s[2]}

		constraints = append(constraints, fk)
	}

	return constraints
}

func (col *ColDesc) getKeyFromDb(db *Database, tableName string) error {
	if col.Name == "" {
		return errors.New("column name is empty")
	}

	result, err := db.Db.Query(`SELECT CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = ?`, db.dbName, tableName, col.Name)
	if err != nil {
		return err
	}

	if result.Next() {
		var constraintName string
		var columnName string
		var referencedTableName sql.NullString
		var referencedColumnName sql.NullString
		err = result.Scan(&constraintName, &columnName, &referencedTableName, &referencedColumnName)
		if err != nil {
			return err
		}

		if strings.ToLower(constraintName) == "primary" {
			col.Key = "primary"
		} else if referencedTableName.Valid && referencedColumnName.Valid {
			col.Key = fmt.Sprintf("foreign,%s,%s,%s", referencedTableName.String, referencedColumnName.String, constraintName)
		}
	}

	return nil
}

func (constraint Constraint) Format(action string) string {
	switch action {
	case "drop":
		return fmt.Sprintf(
			"ALTER TABLE %s DROP FOREIGN KEY %s",
			constraint.TableName,
			constraint.ConstraintName,
		)
	default:
		return fmt.Sprintf(
			"ALTER TABLE %s ADD FOREIGN KEY %s(%s) REFERENCES %s(%s)",
			constraint.TableName,
			constraint.ConstraintName,
			constraint.ColumnName,
			constraint.ReferencedTableName,
			constraint.ReferencedColumnName,
		)
	}
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
		colSql := col.Format(diff.Name)
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

	tableQuery = fmt.Sprintf("ALTER TABLE %s %s;", diff.Name, columns)
	return tableQuery, addConstraintQueries, dropConstraintQueries
}
