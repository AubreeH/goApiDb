package tests

import (
	"github.com/AubreeH/goApiDb/structParsing"
	"testing"
)

func Test_GetTableInfo(t *testing.T) {
	tableInfo, err := structParsing.GetTableInfo(testingEntity1{})
	assertError(t, err)

	assert(t, condition(tableInfo.Name != "test_entity_1", "retrieved table name is incorrect ("+tableInfo.Name+")"))
}
