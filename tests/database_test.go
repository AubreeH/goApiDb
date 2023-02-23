package tests

import (
	"github.com/AubreeH/goApiDb/database"
	"log"
	"testing"
)

func Test_GetTableSqlDescription(t *testing.T) {
	log.Print(database.GetTableSqlDescription[testingEntity1]())
}
