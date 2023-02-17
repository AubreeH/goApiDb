package tests

import (
	"github.com/AubreeH/goApiDb/access"
	"github.com/AubreeH/goApiDb/helpers"
	"testing"
)

type testingEntity struct {
	Id          int64  `json:"id" sql_name:"id" sql_type:"int(64)" sql_key:"PRIMARY" sql_extras:"AUTO_INCREMENT" sql_nullable:"NO" sql_disallow_external_modification:"true"`
	Name        string `json:"name" sql_name:"name" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
	Description string `json:"description" sql_name:"description" sql_type:"VARCHAR(256)" sql_nullable:"NO"`
}

func setupGetById() (testingEntity, error) {
	testEntityName := randSeq(20)
	testEntityDescription := randSeq(20)

	tableName := helpers.GetTableName(testingEntity{})
	id, err := seedTableWithValueInMiddle(
		1000,
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
		return testingEntity{}, nil
	}

	testEntity := testingEntity{
		Id:          id,
		Name:        testEntityName,
		Description: testEntityDescription,
	}
	return testEntity, nil
}

func Test_GetById_Success(t *testing.T) {
	closeFunc, err := setupTable(testingEntity{})
	if err != nil {
		t.Error(err)
		return
	}
	defer closeFunc()

	testEntity, err := setupGetById()
	if err != nil {
		t.Error(err)
		return
	}

	entity, err := access.GetById(db, testingEntity{}, testEntity.Id)
	if err != nil {
		t.Error(err)
		return
	}

	if testEntity.Id != entity.Id {
		t.Error("Ids do not match")
	}

	if testEntity.Name != entity.Name {
		t.Error("Names do not match")
	}

	if testEntity.Description != entity.Description {
		t.Error("Descriptions do not match")
	}
}

func Test_GetById_InvalidId(t *testing.T) {
	closeFunc, err := setupTable(testingEntity{})
	if err != nil {
		t.Error(err)
		return
	}
	defer closeFunc()

	_, err = setupGetById()
	if err != nil {
		t.Error(err)
		return
	}

	entity, err := access.GetById(db, testingEntity{}, "abc")
	if err == nil {
		t.Error("No result should have returned error")
		return
	} else if err.Error() != "unable to find value" {
		t.Error(err)
		return
	}

	if entity.Id != 0 {
		t.Error("Testing entity id should be 0")
	}

	if entity.Name != "" {
		t.Error("Testing entity name should be empty string")
	}

	if entity.Description != "" {
		t.Error("Testing entity description should be empty string")
	}
}
