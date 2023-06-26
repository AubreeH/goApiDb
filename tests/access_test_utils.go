package tests

import (
	"database/sql"
	"time"

	"github.com/AubreeH/goApiDb/entities"
	"github.com/AubreeH/goApiDb/structParsing"
)

type testingEntity1 struct {
	entities.EntityBase `table_name:"test_entity_1"`
	Id                  int64         `json:"id" db_type:"int(64)" db_key:"PRIMARY" db_extras:"AUTO_INCREMENT"`
	Name                string        `json:"name" db_type:"VARCHAR(256)"`
	Description         string        `json:"description" db_type:"VARCHAR(256)"`
	TestEntity2Id       sql.NullInt64 `json:"test_entity_2_id" db_type:"int(64)" db_nullable:"true" db_key:"foreign,test_entity_2,id" parse_struct:"false"`
	Drop                time.Time     `json:"drop" db_type:"datetime" db_nullable:"false" parse_struct:"false"`
	entities.Dates
}

type testingEntity2 struct {
	entities.EntityBase `table_name:"test_entity_2"`
	Id                  int64         `json:"id" db_type:"int(64)" db_key:"PRIMARY" db_extras:"AUTO_INCREMENT" db_disallow_external_modification:"true"`
	Name                string        `json:"name" db_type:"VARCHAR(256)" db_nullable:"NO"`
	Description         string        `json:"description" db_type:"VARCHAR(256)" db_nullable:"NO"`
	TestEntity3Id       sql.NullInt64 `json:"test_entity_3_id" db_type:"int(64)" db_nullable:"true" db_key:"foreign,test_entity_3,id" parse_struct:"false"`
	entities.Dates
}

type testingEntity3 struct {
	entities.EntityBase `table_name:"test_entity_3"`
	Id                  int64  `json:"id" db_type:"int(64)" db_key:"PRIMARY" db_extras:"AUTO_INCREMENT" db_disallow_external_modification:"true"`
	Name                string `json:"name" db_type:"VARCHAR(256)" db_nullable:"NO"`
	Description         string `json:"description" db_type:"VARCHAR(256)" db_nullable:"NO"`
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
		10000,
		tableInfo.Name,
		map[string]string{
			"name":        "string",
			"description": "string",
			"created_at":  "time",
			"updated_at":  "time",
			"drop":        "time",
		},
		map[string]any{
			"name":        testEntityName,
			"description": testEntityDescription,
			"created_at":  testEntityCreatedAt,
			"updated_at":  testEntityUpdatedAt,
			"drop":        time.Now(),
		},
	)
	if err != nil {
		return testingEntity1{}, err
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
		"drop":        "time",
	})
}
