package tests

import (
	"github.com/AubreeH/goApiDb/query"
	"testing"
	"time"
)

func Test_QueryBuilder_Success(t *testing.T) {
	InitDb()

	_, err := setupTables(testingEntity1{}, testingEntity2{}, testingEntity3{})
	//defer closeFunc()

	for i := 0; i < 5000; i++ {
		id, _, err := seedTableWithValueInMiddle(50, "test_entity_3", map[string]string{
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

		if i%2 == 0 {
			id, _, err := seedTableWithValueInMiddle(60, "test_entity_2", map[string]string{
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

			if i%4 == 0 {
				_, _, err := seedTableWithValueInMiddle(70, "test_entity_1", map[string]string{
					"name":        "string",
					"description": "string",
					"created_at":  "time",
					"updated_at":  "time",
				}, map[string]any{
					"name":            randSeq(20),
					"description":     randSeq(20),
					"test_entity2_id": id,
					"created_at":      time.Now(),
					"updated_at":      time.Now(),
				})
				assertError(t, err)
			}
		}
	}

	q := query.NewSelectQuery()
	q.Select("te2.name as Te2name", "te1.id as Id", "te1.name as Te1name", "te3.name as Te3name").
		From(testingEntity1{}, "te1").
		LeftJoin(testingEntity2{}, "te2", "te1.test_entity2_id = te2.id").
		LeftJoin(testingEntity3{}, "te3", "te2.test_entity3_id = te3.id").
		Where("te2.name IS NOT NULL").
		Where("te3.name IS NOT NULL")

	start := time.Now()

	results, err := query.ExecuteQuery(db, q, struct {
		Te2name string
		Id      int64
		Te1name string
		Te3name string
	}{})

	end := time.Now()

	t.Log("QueryBuilder Exec Duration", end.UnixMicro()-start.UnixMicro(), "With", len(results), "results")
	assertError(t, err)
}
