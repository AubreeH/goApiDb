package database

import (
	"errors"
	"github.com/AubreeH/goApiDb/helpers"
)

func getDescriptionDifferences(entityDesc TablDesc, dbDesc TablDesc) (TablDescDiff, error) {
	diff := TablDescDiff{}

	diff.Name = entityDesc.Name

	columnsToAdd, columnsToModify, columnsToDrop, err := getColumnDifferences(entityDesc.Name, entityDesc.Columns, dbDesc.Columns)
	if err != nil {
		return TablDescDiff{}, err
	}

	diff.ColumnsToAdd = columnsToAdd
	diff.ColumnsToModify = columnsToModify
	diff.ColumnsToDrop = columnsToDrop

	constraintsToAdd, constraintsToDrop := getConstraintDifferences(entityDesc.Constraints, dbDesc.Constraints)

	diff.ConstraintsToAdd = constraintsToAdd
	diff.ConstraintsToDrop = constraintsToDrop

	return diff, nil
}

func getColumnDifferences(tableName string, entityColumns []ColDesc, dbColumns []ColDesc) (add []ColDesc, modify []ColDesc, drop []ColDesc, err error) {
	add = []ColDesc{}
	modify = []ColDesc{}
	drop = []ColDesc{}

	for _, entityCol := range entityColumns {
		if entityCol.Name == "" {
			return nil, nil, nil, errors.New("column with empty name provided with entity description")
		}

		dbCol, ok := helpers.ArrFindFunc(dbColumns, func(dbCol ColDesc) bool {
			return entityCol.Name == dbCol.Name
		})
		if !ok {
			add = append(add, entityCol)
		} else {
			entityColSql := entityCol.Format(tableName)
			dbColSql := dbCol.Format(tableName)

			if entityColSql != dbColSql {
				modify = append(modify, entityCol)
			}
		}
	}

	for _, dbCol := range dbColumns {
		if dbCol.Name == "" {
			return nil, nil, nil, errors.New("column with empty name provided with db description")
		}

		_, ok := helpers.ArrFindFunc(entityColumns, func(entityCol ColDesc) bool {
			return entityCol.Name == dbCol.Name
		})

		if !ok {
			drop = append(drop, dbCol)
		}
	}

	return add, modify, drop, nil
}

func getConstraintDifferences(entityConstraints []Constraint, dbConstraints []Constraint) (add []Constraint, drop []Constraint) {
	add = []Constraint{}
	drop = []Constraint{}

	for _, entityConstraint := range entityConstraints {
		_, ok := helpers.ArrFindFunc(dbConstraints, func(dbConstraint Constraint) bool {
			return entityConstraint.TableName == dbConstraint.TableName &&
				entityConstraint.ColumnName == dbConstraint.ColumnName &&
				entityConstraint.ReferencedTableName == dbConstraint.ReferencedTableName &&
				entityConstraint.ReferencedColumnName == dbConstraint.ReferencedColumnName
		})

		if !ok {
			add = append(add, entityConstraint)
		}
	}

	for _, dbConstraint := range dbConstraints {
		_, ok := helpers.ArrFindFunc(entityConstraints, func(entityConstraint Constraint) bool {
			return entityConstraint.TableName == dbConstraint.TableName &&
				entityConstraint.ColumnName == dbConstraint.ColumnName &&
				entityConstraint.ReferencedTableName == dbConstraint.ReferencedTableName &&
				entityConstraint.ReferencedColumnName == dbConstraint.ReferencedColumnName
		})

		if !ok {
			drop = append(drop, dbConstraint)
		}
	}

	return add, drop
}
