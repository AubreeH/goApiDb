package database

import (
	"errors"
	"github.com/AubreeH/goApiDb/helpers"
	"github.com/AubreeH/goApiDb/structParsing"
)

func getDescriptionDifferences(entityDesc structParsing.TablDesc, dbDesc structParsing.TablDesc) (TablDescDiff, error) {
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

func getColumnDifferences(tableName string, entityColumns []structParsing.ColDesc, dbColumns []structParsing.ColDesc) (add []structParsing.ColDesc, modify []structParsing.ColDesc, drop []structParsing.ColDesc, err error) {
	add = []structParsing.ColDesc{}
	modify = []structParsing.ColDesc{}
	drop = []structParsing.ColDesc{}

	for _, entityCol := range entityColumns {
		if entityCol.Name == "" {
			return nil, nil, nil, errors.New("column with empty name provided with entity description")
		}

		dbCol, ok := helpers.ArrFindFunc(dbColumns, func(dbCol structParsing.ColDesc) bool {
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

		_, ok := helpers.ArrFindFunc(entityColumns, func(entityCol structParsing.ColDesc) bool {
			return entityCol.Name == dbCol.Name
		})

		if !ok {
			drop = append(drop, dbCol)
		}
	}

	return add, modify, drop, nil
}

func getConstraintDifferences(entityConstraints []structParsing.Constraint, dbConstraints []structParsing.Constraint) (add []structParsing.Constraint, drop []structParsing.Constraint) {
	add = []structParsing.Constraint{}
	drop = []structParsing.Constraint{}

	for _, entityConstraint := range entityConstraints {
		_, ok := helpers.ArrFindFunc(dbConstraints, func(dbConstraint structParsing.Constraint) bool {
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
		_, ok := helpers.ArrFindFunc(entityConstraints, func(entityConstraint structParsing.Constraint) bool {
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
