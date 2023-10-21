package tests

import "testing"

func setupQueryBuilder(t *testing.T, multiplier int) (map[int]map[string]any, map[int]map[string]any, map[int]map[string]any) {
	InitDb(t)
	setupTables(t, testingEntity1{}, testingEntity2{}, testingEntity3{})
	_, testEntity3SeededValues := seedTable(t, 100*multiplier, "test_entity_3",
		map[string]string{
			"id":          "id",
			"name":        "string",
			"description": "string",
		},
	)

	testEntity2DataValues := make([]map[string]any, 0)
	count := 0
	for k := range testEntity3SeededValues {
		if k%2 == 0 {
			count++
			testEntity2DataValues = append(testEntity2DataValues, map[string]any{
				"id":              nil,
				"name":            randSeq(20),
				"description":     randSeq(20),
				"test_entity3_id": k,
			})
		}
	}

	_, testEntity2SeededValues := seedTable(t, 100*multiplier, "test_entity_2", map[string]string{
		"id":              "id",
		"name":            "string",
		"description":     "string",
		"test_entity3_id": "null",
	}, testEntity2DataValues...)

	testEntity1DataValues := make([]map[string]any, 0)
	count = 0
	for k := range testEntity2SeededValues {
		if k%2 == 0 {
			count++
			testEntity1DataValues = append(testEntity1DataValues, map[string]any{
				"id":              nil,
				"name":            randSeq(20),
				"description":     randSeq(20),
				"test_entity2_id": k,
			})
		}
	}

	_, testEntity1SeededValues := seedTable(t, 100*multiplier, "test_entity_1", map[string]string{
		"id":              "id",
		"name":            "string",
		"description":     "string",
		"test_entity2_id": "null",
	}, testEntity1DataValues...)

	return testEntity3SeededValues, testEntity2SeededValues, testEntity1SeededValues
}
