package tests

import (
	"fmt"
	"testing"

	"github.com/AubreeH/goApiDb/access"
	_ "github.com/go-sql-driver/mysql"
)

func Test_GetById_Success(t *testing.T) {
	testEntity := setupGetById(t)

	entity, err := access.GetById(db, testingEntity1{}, testEntity.Id)

	assert(t,
		e(err),

		condition(testEntity.Id != entity.Id, "ids do no match"),
		condition(testEntity.Name != entity.Name, "names do not match"),
		condition(testEntity.Description != entity.Description, "descriptions do not match"),
	)
}

func Test_GetById_InvalidId(t *testing.T) {
	setupGetById(t)

	entity, err := access.GetById(db, testingEntity1{}, "abc")

	assert(t,
		condition(err != nil && err.Error() != "unable to find value", err.Error()),
		condition(err == nil, "no result should have returned error"),

		condition(entity.Id != 0, "testing entity id should be 0"),
		condition(entity.Name != "", "testing entity name should be empty string"),
		condition(entity.Description != "", "testing entity description should be empty string"),
	)
}

func Test_GetAll_Success(t *testing.T) {
	seededValues := setupGetAll(t)

	results, err := access.GetAll(db, testingEntity1{}, 0)
	assertError(t, err)

	assert(t, condition(
		len(results) != len(seededValues),
		"length of GetAll result differs from length of seeded values (Length of Results:",
		len(results),
		", Length of Seeded Values: ",
		len(seededValues),
		")",
	))

	for _, v := range results {
		seededValue := seededValues[v.Id]

		assert(t,
			condition(seededValue == nil, "unexpected id from result (", v.Id, ")"),
			condition(v.Name != seededValue["name"], "names do not match (", v.Name, "!=", seededValue["name"], ")"),
			condition(v.Description != seededValue["description"], "descriptions do not match (", v.Description, "!=", seededValue["name"], ")"),
		)
	}
}

func Test_Delete_Success(t *testing.T) {
	testEntity := setupGetById(t)

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

func Test_Create_Success(t *testing.T) {
	InitDb(t)
	setupTables(t, false, testingEntity3{})

	testEntityName := randSeq(20)
	testEntityDescription := randSeq(20)

	testEntity := testingEntity3{
		Name:        testEntityName,
		Description: testEntityDescription,
	}

	e1, err := access.Create(db, testEntity)
	assertError(t, err)

	var e2 testingEntity3
	rows, err := db.Db.Query("SELECT id, name, description FROM test_entity_3 WHERE id = ?", e1.Id)
	assertError(t, err)
	for rows.Next() {
		err := rows.Scan(&e2.Id, &e2.Name, &e2.Description)
		assertError(t, err)
	}

	assert(t,
		e(err),

		condition(e1.Id == 0, "id is set"),
		condition(e1.Name == "", "name is set"),
		condition(e1.Description == "", "description is set"),

		condition(e2.Id != e1.Id, "ids do not match"),
		condition(e2.Name != testEntityName, "name does not match"),
		condition(e2.Description != testEntityDescription, "description does not match"),

		condition(e1.Name != testEntityName, "name does not match"),
		condition(e1.Description != testEntityDescription, "description does not match"),
	)
}
