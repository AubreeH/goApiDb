package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/AubreeH/goApiDb/access"
	"github.com/AubreeH/goApiDb/database"
	_ "github.com/go-sql-driver/mysql"
)

func Test_GetById_WithDbObject_Success(t *testing.T) {
	testEntity := setupGetById(t)

	entity, err := access.GetById(db, testingEntity1{}, testEntity.Id)

	assert(t,
		e(err),

		condition(testEntity.Id != entity.Id, "ids do no match"),
		condition(testEntity.Name != entity.Name, "names do not match"),
		condition(testEntity.Description != entity.Description, "descriptions do not match"),
	)
}

func Test_GetById_WithDbObject_InvalidId(t *testing.T) {
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

func Test_GetAll_WithDbObject_Success(t *testing.T) {
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

func Test_Delete_WithDbObject_Success(t *testing.T) {
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

func Test_Delete_WithTransactionObject_Success(t *testing.T) {
	// TODO: Setup Delete test for transaction
}

func Test_Create_WithDbObject_Success(t *testing.T) {
	InitDb(t)
	setupTables(t, true, testingEntity3{})

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

func Test_Create_WithTransactionObject(t *testing.T) {
	InitDb(t)
	setupTables(t, true, testingEntity2{}, testingEntity3{})

	testEntity1Name := randSeq(20)
	testEntity1Description := randSeq(20)
	testEntity2Name := randSeq(20)
	testEntity2Description := randSeq(20)

	err := db.Transaction(func(tx *database.Transaction) error {
		innerEntity1, err := access.Create(tx, testingEntity3{
			Name:        testEntity1Name,
			Description: testEntity1Description,
		})
		if err != nil {
			return err
		} else if innerEntity1.Id == 0 {
			return fmt.Errorf("e1 id is not set")
		}

		innerEntity2, err := access.Create(tx, testingEntity2{
			Name:          testEntity2Name,
			Description:   testEntity2Description,
			TestEntity3Id: &innerEntity1.Id,
		})
		if err != nil {
			return err
		} else if innerEntity2.Id == 0 {
			return fmt.Errorf("e2 id is not set")
		}

		return nil
	})
	assertError(t, err)

	outerEntity1 := testingEntity3{}
	rows, err := db.Db.Query("SELECT id, name, description FROM test_entity_3 WHERE name = ?", testEntity1Name)
	assertError(t, err)
	if !rows.Next() {
		assertError(t, fmt.Errorf("no rows returned when querying for first test entity"))
	}

	err = rows.Scan(&outerEntity1.Id, &outerEntity1.Name, &outerEntity1.Description)
	assertError(t, err)

	outerEntity2 := testingEntity2{}
	rows, err = db.Db.Query("SELECT id, name, description, test_entity3_id FROM test_entity_2 WHERE name = ?", testEntity2Name)
	assertError(t, err)
	if !rows.Next() {
		assertError(t, fmt.Errorf("no rows returned when querying for second test entity"))
	}

	err = rows.Scan(&outerEntity2.Id, &outerEntity2.Name, &outerEntity2.Description, &outerEntity2.TestEntity3Id)
	assertError(t, err)

	assert(t,
		condition(outerEntity1.Id == 0, "first entity id is not set"),
		condition(outerEntity1.Name == "", "first entity name is not set"),
		condition(outerEntity1.Description == "", "first entity description is not set"),
		condition(outerEntity1.Name != testEntity1Name, "first entity name does not match"),
		condition(outerEntity1.Description != testEntity1Description, "first entity description does not match"),
		condition(outerEntity2.Id == 0, "second entity id is not set"),
		condition(outerEntity2.Name == "", "second entity name is not set"),
		condition(outerEntity2.Description == "", "second entity description is not set"),
		condition(outerEntity2.Name != testEntity2Name, "second entity name does not match"),
		condition(outerEntity2.Description != testEntity2Description, "second entity description does not match"),
		condition(outerEntity1.Id != *outerEntity2.TestEntity3Id, "first entity and second entity ids do not match"),
	)

	testEntity3Name := randSeq(20)
	testEntity3Description := randSeq(20)
	testEntity4Name := randSeq(20)
	testEntity4Description := randSeq(20)

	err = db.Transaction(func(tx *database.Transaction) error {
		innerEntity3, err := access.Create(tx, testingEntity3{
			Name:        testEntity3Name,
			Description: testEntity3Description,
		})
		if err != nil {
			return err
		} else if innerEntity3.Id == 0 {
			return fmt.Errorf("e2 id is not set")
		}

		innerEntity4, err := access.Create(tx, testingEntity2{
			Name:          testEntity4Name,
			Description:   testEntity4Description,
			TestEntity3Id: &innerEntity3.Id,
		})
		if err != nil {
			return err
		} else if innerEntity4.Id == 0 {
			return fmt.Errorf("e2 id is not set")
		}

		return errors.New("rollback")
	})
	if err != nil && err.Error() != "rollback" {
		assertError(t, err)
	}

	outerEntity3 := testingEntity3{}
	rows, err = db.Db.Query("SELECT id, name, description FROM test_entity_3 WHERE name = ?", testEntity3Name)
	assertError(t, err)
	for rows.Next() {
		err := rows.Scan(&outerEntity3.Id, &outerEntity3.Name, &outerEntity3.Description)
		assertError(t, err)
	}

	outerEntity4 := testingEntity2{}
	rows, err = db.Db.Query("SELECT id, name, description, test_entity3_id FROM test_entity_2 WHERE name = ?", testEntity4Name)
	assertError(t, err)
	for rows.Next() {
		err := rows.Scan(&outerEntity4.Id, &outerEntity4.Name, &outerEntity4.Description, &outerEntity4.TestEntity3Id)
		assertError(t, err)
	}

	assert(t,
		condition(outerEntity3.Id != 0, "e3 id is set"),
		condition(outerEntity3.Name != "", "e3 name is set"),
		condition(outerEntity3.Description != "", "e3 description is set"),
		condition(outerEntity4.Id != 0, "e4 id is set"),
		condition(outerEntity4.Name != "", "e4 name is set"),
		condition(outerEntity4.Description != "", "e4 description is set"),
		condition(outerEntity4.TestEntity3Id != nil, "e4 test entity 3 id is set"),
	)

}
