package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/AubreeH/goApiDb/access"
	_ "github.com/go-sql-driver/mysql"
)

func Test_GetById_Success(t *testing.T) {
	InitDb()

	closeFunc, err := setupTables(testingEntity1{}, testingEntity2{}, testingEntity3{})
	assert(t, condition(err != nil, err))
	defer closeFunc()

	testEntity, err := setupGetById()
	assertError(t, err)

	entity, timedResult, err := access.GetByIdTimed(db, testingEntity1{}, testEntity.Id)
	assertError(t, err)

	t.Log("Overall Duration", fmt.Sprint(timedResult.OverallDuration, "µs"), "At Id of", testEntity.Id)
	t.Log("Query Build Duration", fmt.Sprint(timedResult.BuildQueryDuration, "µs"))
	t.Log("Query Exec Duration", fmt.Sprint(timedResult.QueryExecDuration, "µs"))
	t.Log("Format Result Duration", fmt.Sprint(timedResult.FormatResultDuration, "µs"))

	assert(t,
		condition(testEntity.Id != entity.Id, "ids do no match"),
		condition(testEntity.Name != entity.Name, "names do not match"),
		condition(testEntity.Description != entity.Description, "descriptions do not match"),
	)
}

func Test_GetById_InvalidId(t *testing.T) {
	InitDb()

	closeFunc, err := setupTables(testingEntity1{}, testingEntity2{}, testingEntity3{})
	assert(t, condition(err != nil, err))
	defer closeFunc()

	_, err = setupGetById()
	assert(t, condition(err != nil, err))

	start := time.Now()
	entity, err := access.GetById(db, testingEntity1{}, "abc")
	end := time.Now()
	assert(t, condition(err != nil && err.Error() != "unable to find value", err.Error()))

	t.Log("GetById Duration", fmt.Sprint(end.UnixMicro()-start.UnixMicro(), "µs"))

	assert(t,
		condition(entity.Id != 0, "testing entity id should be 0"),
		condition(entity.Name != "", "testing entity name should be empty string"),
		condition(entity.Description != "", "testing entity description should be empty string"),
	)
	assert(t, condition(err == nil, "no result should have returned error"))
}

func Test_GetAll_Success(t *testing.T) {
	InitDb()

	closeFunc, err := setupTables(testingEntity1{}, testingEntity2{}, testingEntity3{})
	assert(t, condition(err != nil, err))
	defer closeFunc()

	seededValues, err := setupGetAll()
	assert(t, condition(err != nil, err))

	results, timedResult, err := access.GetAllTimed(db, testingEntity1{}, 0)
	assert(t, condition(err != nil, err))

	t.Log("Overall Duration", fmt.Sprint(timedResult.OverallDuration, "µs"))
	t.Log("Query Build Duration", fmt.Sprint(timedResult.BuildQueryDuration, "µs"))
	t.Log("Query Exec Duration", fmt.Sprint(timedResult.QueryExecDuration, "µs"))
	t.Log("Format Result Duration", fmt.Sprint(timedResult.FormatResultDuration, "µs"))

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

func Test_Delete_Success(t *testing.T) {
	InitDb()

	closeFunc, err := setupTables(testingEntity1{}, testingEntity2{}, testingEntity3{})
	assert(t, condition(err != nil, err))
	defer closeFunc()

	testEntity, err := setupGetById()
	assertError(t, err)

	timedResult, err := access.DeleteTimed(db, testingEntity1{}, testEntity.Id)
	assertError(t, err)

	t.Log("Overall Duration", fmt.Sprint(timedResult.OverallDuration, "µs"), "At Id of", testEntity.Id)
	t.Log("Query Build Duration", fmt.Sprint(timedResult.BuildQueryDuration, "µs"))
	t.Log("Query Exec Duration", fmt.Sprint(timedResult.QueryExecDuration, "µs"))

	entity, err := access.GetById(db, testingEntity1{}, testEntity.Id)

	assert(t,
		condition(err.Error() != "unable to find value", "error does not match", err),
		condition(entity.Id != 0, "id is set"),
		condition(entity.Name != "", "name is set"),
		condition(entity.Description != "", "description is set"),
	)
}
