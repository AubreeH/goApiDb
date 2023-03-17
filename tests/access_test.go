package tests

import (
	"github.com/AubreeH/goApiDb/access"
	"testing"
	"time"
)

func Test_GetById_Success(t *testing.T) {
	InitDb()

	_, err := setupTables(testingEntity1{}, testingEntity2{})
	assert(t, condition(err != nil, err))
	//defer closeFunc()

	testEntity, err := setupGetById()
	assert(t, condition(err != nil, err))

	start := time.Now()
	entity, err := access.GetById(db, testingEntity1{}, testEntity.Id)
	end := time.Now()
	assert(t, condition(err != nil, err))

	t.Log("GetById Duration", end.UnixMicro()-start.UnixMicro(), "At Id of", testEntity.Id)

	assert(t,
		condition(testEntity.Id != entity.Id, "ids do no match"),
		condition(testEntity.Name != entity.Name, "names do not match"),
		condition(testEntity.Description != entity.Description, "descriptions do not match"),
	)
}

func Test_GetById_InvalidId(t *testing.T) {
	InitDb()

	closeFunc, err := setupTables(testingEntity1{}, testingEntity2{})
	assert(t, condition(err != nil, err))
	defer closeFunc()

	_, err = setupGetById()
	assert(t, condition(err != nil, err))

	start := time.Now()
	entity, err := access.GetById(db, testingEntity1{}, "abc")
	end := time.Now()
	assert(t, condition(err != nil && err.Error() != "unable to find value", err.Error()))

	t.Log("GetById Duration", end.UnixMicro()-start.UnixMicro())

	assert(t,
		condition(entity.Id != 0, "testing entity id should be 0"),
		condition(entity.Name != "", "testing entity name should be empty string"),
		condition(entity.Description != "", "testing entity description should be empty string"),
	)
	assert(t, condition(err == nil, "no result should have returned error"))
}

func Test_GetAll_Success(t *testing.T) {
	InitDb()

	_, err := setupTables(testingEntity1{}, testingEntity2{})
	assert(t, condition(err != nil, err))
	//defer closeFunc()

	seededValues, err := setupGetAll()
	assert(t, condition(err != nil, err))

	start := time.Now()
	results, err := access.GetAll(db, testingEntity1{}, 0)
	end := time.Now()
	assert(t, condition(err != nil, err))

	t.Log("GetAll Duration", end.UnixMicro()-start.UnixMicro())

	assert(t, condition(len(results) != len(seededValues), "length of GetAll result differs from length of seeded values (Length of Results:", len(results), ", Length of Seeded Values: ", len(seededValues), ")"))

	for _, v := range results {
		seededValue := seededValues[v.Id]

		assert(t, condition(seededValue == nil, "unexpected id from result (", v.Id, ")"))

		assert(t,
			condition(v.Name != seededValue["name"], "names do not match (", v.Name, "!=", seededValue["name"], ")"),
			condition(v.Description != seededValue["description"], "descriptions do not match (", v.Description, "!=", seededValue["name"], ")"),
		)
	}
}
