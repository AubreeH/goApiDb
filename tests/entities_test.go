package tests

import (
	"github.com/AubreeH/goApiDb/entities"
	"testing"
)

func Test_GetTableInfo(t *testing.T) {
	tableInfo, err := entities.GetTableInfo(testingEntity1{})
	assertError(t, err)

	assert(t, condition(tableInfo.Name != "test_entity_1", "retrieved table name is incorrect ("+tableInfo.Name+")"))
}
