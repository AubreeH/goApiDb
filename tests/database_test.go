package tests

import (
	"github.com/AubreeH/goApiDb/database"
	"log"
	"testing"
)

func Test_GetTableSqlDescription(t *testing.T) {
	InitDb()

	log.Print(database.GetUpdateTableQueries[testingEntity1](db))
}
