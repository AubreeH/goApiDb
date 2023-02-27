package tests

import (
	"github.com/AubreeH/goApiDb/database"
	"log"
	"testing"
)

func Test_GetTableSqlDescription(t *testing.T) {
	InitDb()

	entityTableDescription, err := database.GetTableSqlDescriptionFromEntity[testingEntity1]()
	assertError(t, err)

	sqlTableDescription, err := database.GetTableSqlDescriptionFromDb(db, entityTableDescription.Name)
	assertError(t, err)

	log.Print(sqlTableDescription)
}
