package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/AubreeH/goApiDb/query"
)

func Test_QueryBuilder_Success(t *testing.T) {
	InitDb()

	fmt.Println("Setting up tables")
	closeFunc, err := setupTables(testingEntity1{}, testingEntity2{}, testingEntity3{})
	assertError(t, err)
	defer closeFunc()
	fmt.Println("Setting up tables - Done")

	fmt.Println("Seeding tables")
	testEntity1Count := 0
	testEntity2Count := 0
	testEntity3Count := 0
	for i := 0; i < 5000; i++ {
		id, _, err := seedTableWithValueInMiddle(500, "test_entity_3", map[string]string{
			"name":        "string",
			"description": "string",
			"created_at":  "time",
			"updated_at":  "time",
		}, map[string]any{
			"name":        randSeq(20),
			"description": randSeq(20),
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		})
		assertError(t, err)

		testEntity3Count += 500

		if i%2 == 0 {
			id, _, err := seedTableWithValueInMiddle(600, "test_entity_2", map[string]string{
				"name":        "string",
				"description": "string",
				"created_at":  "time",
				"updated_at":  "time",
			}, map[string]any{
				"name":            randSeq(20),
				"description":     randSeq(20),
				"test_entity3_id": id,
				"created_at":      time.Now(),
				"updated_at":      time.Now(),
			})
			assertError(t, err)

			testEntity2Count += 600

			if i%4 == 0 {
				_, _, err := seedTableWithValueInMiddle(700, "test_entity_1", map[string]string{
					"name":        "string",
					"description": "string",
					"created_at":  "time",
					"updated_at":  "time",
					"drop":        "time",
				}, map[string]any{
					"name":            randSeq(20),
					"description":     randSeq(20),
					"test_entity2_id": id,
					"created_at":      time.Now(),
					"updated_at":      time.Now(),
					"drop":            time.Now(),
				})
				assertError(t, err)

				testEntity1Count += 700
			}
		}
	}
	fmt.Printf("Seeding tables - Done (%d test_entity_3, %d test_entity_2, %d test_entity_1)\n", testEntity3Count, testEntity2Count, testEntity1Count)

	fmt.Println("Querying")
	start := time.Now()

	q := query.NewSelectQuery()
	q.Select("te2.name as Te2name", "te1.id as Id", "te1.name as Te1name", "te3.name as Te3name").
		From(testingEntity1{}, "te1").
		LeftJoin(testingEntity2{}, "te2", "te1.test_entity2_id = te2.id").
		LeftJoin(testingEntity3{}, "te3", "te2.test_entity3_id = te3.id").
		Where("te2.name IS NOT NULL").
		Where("te3.name IS NOT NULL")

	results, err := query.ExecuteQuery(db, q, struct {
		Te2name string
		Id      int64
		Te1name string
		Te3name string
	}{})

	end := time.Now()
	fmt.Println("Querying - Done")

	duration := end.UnixMicro() - start.UnixMicro()
	t.Log("QueryBuilder Exec Duration", fmt.Sprint(duration, "µs"), "With", len(results), fmt.Sprintf("results (Average: %dµs)", duration/int64(len(results))))
	assertError(t, err)
}

func Test_QueryBuilderPaginated_Success(t *testing.T) {
	InitDb()

	fmt.Println("Setting up table")
	closeFunc, err := setupTables(testingEntity1{}, testingEntity2{}, testingEntity3{})
	assertError(t, err)
	defer closeFunc()
	fmt.Println("Setting up table - Done")

	fmt.Println("Seeding table")
	_, err = seedTable(500, "test_entity_1", map[string]string{
		"name":        "string",
		"description": "string",
		"created_at":  "time",
		"updated_at":  "time",
		"drop":        "time",
	})
	assertError(t, err)
	fmt.Println("Seeding table - Done")

	q := query.NewSelectQuery()
	q.Select("te1.id", "te1.name")
	q.From(testingEntity1{}, "te1")
	allResults, err := query.ExecuteQuery(db, q, struct {
		Id   int64
		Name string
	}{})
	assertError(t, err)

	q.Paginated(25, 0)
	results1, err := query.ExecuteQuery(db, q, struct {
		Id   int64
		Name string
	}{})
	assertError(t, err)

	q.Paginated(25, 1)
	results2, err := query.ExecuteQuery(db, q, struct {
		Id   int64
		Name string
	}{})
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
