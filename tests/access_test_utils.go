package tests

import (
	"github.com/AubreeH/goApiDb/helpers"
)

type testingEntity1 struct {
	Id            int64  `json:"id" sql_name:"id" sql_type:"int(64)" sql_key:"PRIMARY" sql_extras:"AUTO_INCREMENT" sql_nullable:"NO" sql_disallow_external_modification:"true"`
	Name          string `json:"name" sql_name:"name" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	Description   string `json:"description" sql_name:"description" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	TestEntity2Id int64  `json:"test_entity_2_id" sql_name:"test_entity_2_id" sql_key:"foreign,testing_entity2s,id"`
}

type testingEntity2 struct {
	Id          int64  `json:"id" sql_name:"id" sql_type:"int(64)" sql_key:"PRIMARY" sql_extras:"AUTO_INCREMENT" sql_nullable:"NO" sql_disallow_external_modification:"true"`
	Name        string `json:"name" sql_name:"name" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	Description string `json:"description" sql_name:"description" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
}

func setupGetById() (output testingEntity1, err error) {
	testEntityName := randSeq(20)
	testEntityDescription := randSeq(20)

	tableName := helpers.GetTableName(testingEntity1{})
	id, _, err := seedTableWithValueInMiddle(
		10000,
		tableName,
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
	tableName := helpers.GetTableName(testingEntity1{})
	return seedTable(10000, tableName, map[string]string{
		"name":        "string",
		"description": "string",
	})
}
