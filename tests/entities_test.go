package tests

import (
	"fmt"
	"testing"

	"github.com/AubreeH/goApiDb/structParsing"
)

func Test_GetTableInfo(t *testing.T) {
	tableInfo, err := structParsing.GetTableInfo(testingEntity1{})
	assertError(t, err)

	assert(t,
		condition(tableInfo.Name != "test_entity_1", fmt.Sprintf("Table name does not match (expected: %s, actual: %s)", "test_entity_1", tableInfo.Name)),
		condition(!tableInfo.IsValid, "Table is not valid"),
		condition(tableInfo.PrimaryKey != "id", fmt.Sprintf("Primary key does not match (expected: %s, actual: %s)", "id", tableInfo.PrimaryKey)),
		condition(tableInfo.SoftDeletes != "deleted_at", fmt.Sprintf("Soft deletes does not match (expected: %s, actual: %s)", "deleted_at", tableInfo.SoftDeletes)),
	)
}
