package structParsing

import (
	"fmt"
	"strings"

	"github.com/AubreeH/goApiDb/helpers"
)

type TablDesc struct {
	Name        string
	Columns     []ColDesc
	Constraints []Constraint
}

type ColDesc struct {
	Name                         string
	Type                         string
	Key                          string
	Extras                       string
	Nullable                     string
	Default                      string
	DisallowExternalModification bool
	Pointer                      any
	Value                        any
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

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (%s)", tabl.Name, strings.Join(columns, ", ")), constraints
}

func (col *ColDesc) Format(tableName string, action ...string) string {
	var s []string

	if len(action) > 0 && strings.ToLower(action[0]) == "drop" {
		return col.Name
	}

	helpers.ArrAdd(&s, fmt.Sprintf("`%s`", col.Name), col.Type, FormatKey(col.Key), col.Nullable, col.Default, col.Extras)

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

func (constraint Constraint) Format(action string) string {
	switch action {
	case "drop":
		return fmt.Sprintf(
			"ALTER TABLE `%s` DROP FOREIGN KEY `%s`",
			constraint.TableName,
			constraint.ConstraintName,
		)
	default:
		return fmt.Sprintf(
			"ALTER TABLE `%s` ADD FOREIGN KEY `%s`(`%s`) REFERENCES `%s`(`%s`)",
			constraint.TableName,
			constraint.ConstraintName,
			constraint.ColumnName,
			constraint.ReferencedTableName,
			constraint.ReferencedColumnName,
		)
	}
}
