package tests

import (
	"fmt"
	"testing"

	"github.com/AubreeH/goApiDb/query"
)

func Test_QueryBuilder_Success(t *testing.T) {
	testEntity3DataValues, testEntity2DataValues, testEntity1DataValues := setupQueryBuilder(t, 100)

	q := query.Select(struct {
		Id      int    `s:"te1.id"`
		Te1name string `s:"te1.name"`
		Te2name string `s:"te2.name"`
		Te3name string `s:"te3.name"`
	}{}).
		From(testingEntity1{}, "te1").
		LeftJoin(testingEntity2{}, "te2", "te1.test_entity2_id = te2.id").
		LeftJoin(testingEntity3{}, "te3", "te2.test_entity3_id = te3.id").
		Where(
			"te2.name IS NOT NULL",
			"te3.name IS NOT NULL",
		)

	results, err := q.All(db)
	assertError(t, err)

	for _, v := range results {

		assert(t,
			condition(v.Te1name == "", "te1.name is empty"),
			condition(v.Te2name == "", "te2.name is empty"),
			condition(v.Te3name == "", "te3.name is empty"),
		)

		te1DataValue := testEntity1DataValues[v.Id]
		te2DataValue := testEntity2DataValues[te1DataValue["test_entity2_id"].(int)]
		te3DataValue := testEntity3DataValues[te2DataValue["test_entity3_id"].(int)]

		assert(t,
			condition(te1DataValue["name"] != v.Te1name, "te1.name does not match"),
			condition(te2DataValue["name"] != v.Te2name, "te2.name does not match"),
			condition(te3DataValue["name"] != v.Te3name, "te3.name does not match"),
		)
	}
}

func Test_QueryBuilderPaginated_Success(t *testing.T) {
	InitDb(t)
	setupTables(t, testingEntity1{}, testingEntity2{}, testingEntity3{})

	seedTable(t, 500, "test_entity_1", map[string]string{
		"name":        "string",
		"description": "string",
	})

	q := query.Select(
		struct {
			Id   int64  `s:"te1.id"`
			Name string `s:"te1.name"`
		}{},
	).From(testingEntity1{}, "te1")

	allResults, err := q.All(db)
	assertError(t, err)

	results1, err := q.Paginated(db, 25, 0)
	assertError(t, err)

	results2, err := q.Paginated(db, 25, 1)
	assertError(t, err)

	assert(t,
		condition(len(results1) > 25, "More values retrieved than limit"),
		condition(len(results1) < 25, "Less values retrieved than limit"),
		condition(!checkArraysEqual(results1, allResults[0:25]), "Results do not match", fmt.Sprintln(results1), fmt.Sprintln(allResults[0:25])),
		condition(len(results2) > 25, "More values retrieved than limit"),
		condition(len(results2) < 25, "Less values retrieved than limit"),
		condition(!checkArraysEqual(results2, allResults[25:50]), "Results do not match", fmt.Sprintln(results2), fmt.Sprintln(allResults[0:25])),
	)
}

func checkArraysEqual[T comparable](arrays ...[]T) bool {
	if len(arrays) == 0 || len(arrays) == 1 {
		return true
	}

	firstArray := arrays[0]

	for _, a := range arrays[1:] {
		if len(firstArray) != len(a) {
			return false
		}

		for i, v := range a {
			if v != firstArray[i] {
				return false
			}
		}
	}

	return true
}
