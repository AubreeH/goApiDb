package tests

import (
	"testing"
	"time"

	"github.com/AubreeH/goApiDb/entities"
	"github.com/AubreeH/goApiDb/structParsing"
)

type testingEntity1 struct {
	entities.EntityBase `table_name:"test_entity_1"`
	Id                  int    `json:"id" db_type:"int" db_key:"PRIMARY" db_extras:"AUTO_INCREMENT"`
	Name                string `json:"name" db_type:"VARCHAR(256)"`
	Description         string `json:"description" db_type:"VARCHAR(256)"`
	TestEntity2Id       *int64 `json:"test_entity_2_id" db_type:"int(64)" db_nullable:"true" db_key:"foreign,test_entity_2,id" parse_struct:"false"`
	entities.Dates
}

type testingEntity2 struct {
	entities.EntityBase `table_name:"test_entity_2"`
	Id                  int64  `json:"id" db_type:"int(64)" db_key:"PRIMARY" db_extras:"AUTO_INCREMENT" db_disallow_external_modification:"true"`
	Name                string `json:"name" db_type:"VARCHAR(256)" db_nullable:"NO"`
	Description         string `json:"description" db_type:"VARCHAR(256)" db_nullable:"NO"`
	TestEntity3Id       *int64 `json:"test_entity_3_id" db_type:"int(64)" db_nullable:"true" db_key:"foreign,test_entity_3,id" parse_struct:"false"`
	entities.Dates
}

type testingEntity3 struct {
	entities.EntityBase `table_name:"test_entity_3"`
	Id                  int64  `json:"id" db_type:"int(64)" db_key:"PRIMARY" db_extras:"AUTO_INCREMENT" db_disallow_external_modification:"true"`
	Name                string `json:"name" db_type:"VARCHAR(256)" db_nullable:"NO"`
	Description         string `json:"description" db_type:"VARCHAR(256)" db_nullable:"NO"`
	entities.Dates
}

func setupGetById(t *testing.T) (output testingEntity1) {
	t.Helper()
	InitDb(t)
	setupTables(t, testingEntity1{}, testingEntity2{}, testingEntity3{})
	testEntityName := randSeq(20)
	testEntityDescription := randSeq(20)
	testEntityCreatedAt := time.Now()
	testEntityUpdatedAt := time.Now()

	tableInfo, err := structParsing.GetTableInfo(testingEntity1{})
	assertError(t, err)

	ids, _ := seedTable(t, 10000, tableInfo.Name,
		map[string]string{
			"name":        "string",
			"description": "string",
			"created_at":  "time",
			"updated_at":  "time",
		},
		map[string]any{
			"name":        testEntityName,
			"description": testEntityDescription,
			"created_at":  testEntityCreatedAt,
			"updated_at":  testEntityUpdatedAt,
		},
	)
	assertError(t, err)

	testEntity := testingEntity1{
		Id:          ids[0],
		Name:        testEntityName,
		Description: testEntityDescription,
	}

	testEntity.CreatedAt.CreatedAt.Time = testEntityCreatedAt
	testEntity.UpdatedAt.UpdatedAt.Time = testEntityUpdatedAt
	return testEntity
}

func setupGetAll(t *testing.T) (expectedValue map[int]map[string]any) {
	t.Helper()
	InitDb(t)
	setupTables(t, testingEntity1{}, testingEntity2{}, testingEntity3{})
	tableInfo, err := structParsing.GetTableInfo(testingEntity1{})
	if err != nil {
		return nil
	}
	_, seededValues := seedTable(t, 10000, tableInfo.Name, map[string]string{
		"name":        "string",
		"description": "string",
		"created_at":  "time",
		"updated_at":  "time",
	})
	return seededValues
}
