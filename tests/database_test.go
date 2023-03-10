package tests

import (
	"github.com/AubreeH/goApiDb/database"
	"testing"
)

func Test_GetTableSqlDescription(t *testing.T) {
	InitDb()

	tableSql, constraintsToAdd, constraintsToDrop, err := database.GetUpdateTableQueriesForEntity(db, testingEntity2{})
	assertError(t, err)

	p(tableSql)

	for _, constraint := range constraintsToDrop {
		_, err = db.Db.Exec(constraint)
		assertError(t, err)
	}

	if tableSql != "" {
		_, err = db.Db.Exec(tableSql)
		assertError(t, err)
	}

	p(constraintsToAdd)

	for _, constraint := range constraintsToAdd {
		_, err = db.Db.Exec(constraint)
		assertError(t, err)
	}

	//err = dropTable[testingEntity2]()
	assertError(t, err)
}

func Test_GetUpdateTableQueriesForEntities(t *testing.T) {
	InitDb()

	tableQueries, addConstraintsQueries, dropConstraintsQueries, err := database.GetUpdateTableQueriesForEntities(db, testingEntity1{}, testingEntity2{})
	assertError(t, err)

	for _, query := range dropConstraintsQueries {
		_, err = db.Db.Exec(query)
		assertError(t, err)
	}

	for _, query := range tableQueries {
		_, err = db.Db.Exec(query)
		assertError(t, err)
	}

	for _, query := range addConstraintsQueries {
		_, err = db.Db.Exec(query)
		assertError(t, err)
	}

	err = dropTable[testingEntity1]()
	assertError(t, err)

	err = dropTable[testingEntity2]()
	assertError(t, err)
}

func Test_BuildTables(t *testing.T) {
	InitDb()

	err := database.BuildTables(db, testingEntity1{}, testingEntity2{})
	assertError(t, err)

	err = dropTable[testingEntity1]()
	assertError(t, err)

	err = dropTable[testingEntity2]()
	assertError(t, err)
}

func Test_RunMigrations(t *testing.T) {
	InitDb()

	err := database.RunMigrations(db, "./testMigrations/")
	assertError(t, err)
}
