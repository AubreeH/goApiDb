package tests

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/entities"
)

type testingEntity1 struct {
	entities.EntityBase `table_name:"test_entity_1"`
	Id                  int64         `json:"id" sql_name:"id" sql_type:"int(64)" sql_key:"PRIMARY" sql_extras:"AUTO_INCREMENT" sql_nullable:"NO" sql_disallow_external_modification:"true"`
	Name                string        `json:"name" sql_name:"name" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	Description         string        `json:"description" sql_name:"description" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	TestEntity2Id       sql.NullInt64 `json:"test_entity_2_id" sql_name:"test_entity_2_id" sql_type:"int(64)" sql_nullable:"true" sql_key:"foreign,test_entity_2,id" parse_struct:"false"`
	entities.UpdatedAt
}

type testingEntity2 struct {
	entities.EntityBase `table_name:"test_entity_2"`
	Id                  int64  `json:"id" sql_name:"id" sql_type:"int(64)" sql_key:"PRIMARY" sql_extras:"AUTO_INCREMENT" sql_nullable:"NO" sql_disallow_external_modification:"true"`
	Name                string `json:"name" sql_name:"name" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	Description         string `json:"description" sql_name:"description" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	entities.CreatedAt
	entities.UpdatedAt
}

func setupGetById() (output testingEntity1, err error) {
	testEntityName := randSeq(20)
	testEntityDescription := randSeq(20)

	tableInfo, err := entities.GetTableInfo(testingEntity1{})
	if err != nil {
		return testingEntity1{}, err
	}

	id, _, err := seedTableWithValueInMiddle(
		10000,
		tableInfo.Name,
		map[string]string{
			"name":        "string",
			"description": "string",
		},
		map[string]any{
			"name":        testEntityName,
			"description": testEntityDescription,
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
	return testEntity, nil
}

func setupGetAll() (expectedValue map[int64]map[string]any, err error) {
	tableInfo, err := entities.GetTableInfo(testingEntity1{})
	if err != nil {
		return nil, err
	}
	return seedTable(10000, tableInfo.Name, map[string]string{
		"name":        "string",
		"description": "string",
	})
}
