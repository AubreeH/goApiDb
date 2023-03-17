package tests

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/entities"
	"github.com/AubreeH/goApiDb/structParsing"
	"time"
)

type testingEntity1 struct {
	entities.EntityBase `table_name:"test_entity_1"`
	Id                  int64         `json:"id" sql_type:"int(64)" sql_key:"PRIMARY" sql_extras:"AUTO_INCREMENT"`
	Name                string        `json:"name" sql_type:"VARCHAR(256)"`
	Description         string        `json:"description" sql_type:"VARCHAR(256)"`
	TestEntity2Id       sql.NullInt64 `json:"test_entity_2_id" sql_type:"int(64)" sql_nullable:"true" sql_key:"foreign,test_entity_2,id" parse_struct:"false"`
	entities.Dates
}

type testingEntity2 struct {
	entities.EntityBase `table_name:"test_entity_2"`
	Id                  int64  `json:"id" sql_type:"int(64)" sql_key:"PRIMARY" sql_extras:"AUTO_INCREMENT" sql_disallow_external_modification:"true"`
	Name                string `json:"name" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	Description         string `json:"description" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	entities.Dates
}

func setupGetById() (output testingEntity1, err error) {
	testEntityName := randSeq(20)
	testEntityDescription := randSeq(20)
	testEntityCreatedAt := time.Now()
	testEntityUpdatedAt := time.Now()

	tableInfo, err := structParsing.GetTableInfo(testingEntity1{})
	if err != nil {
		return testingEntity1{}, err
	}

	id, _, err := seedTableWithValueInMiddle(
		1000000,
		tableInfo.Name,
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
	if err != nil {
		return testingEntity1{}, nil
	}

	testEntity := testingEntity1{
		Id:          id,
		Name:        testEntityName,
		Description: testEntityDescription,
	}

	testEntity.CreatedAt.CreatedAt.Time = testEntityCreatedAt
	testEntity.UpdatedAt.UpdatedAt.Time = testEntityUpdatedAt
	return testEntity, nil
}

func setupGetAll() (expectedValue map[int64]map[string]any, err error) {
	tableInfo, err := structParsing.GetTableInfo(testingEntity1{})
	if err != nil {
		return nil, err
	}
	return seedTable(10000, tableInfo.Name, map[string]string{
		"name":        "string",
		"description": "string",
		"created_at":  "time",
		"updated_at":  "time",
	})
}
