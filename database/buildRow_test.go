package database

import (
	"testing"

	"github.com/AubreeH/goApiDb/entities"
)

type testStruct struct {
	entities.EntityBase `table_name:"test_table"`
	Property1           string
	Property2           int
	Property3           float64
	Property4           bool
	Property5           uint
	Property6           uint64
	Property7           uint32
	Property8           uint16
	Property9           uint8
}

func TestGetEntityConstruction(t *testing.T) {
	var test testStruct
	fields, ptr, tableName, err := getEntityConstruction(&test)
	if err != nil {
		t.Errorf("error getting entity construction: %s", err)
		t.FailNow()
	}

	if tableName != "test_table" {
		t.Log(tableName)
		t.Error("table name does not match")
		t.FailNow()
	}

	if len(fields) != 9 {
		t.Error("incorrect number of fields")
		t.FailNow()
	}

	property1 := testMapField[*string](t, fields, "property1")
	property2 := testMapField[*int](t, fields, "property2")
	property3 := testMapField[*float64](t, fields, "property3")
	property4 := testMapField[*bool](t, fields, "property4")
	property5 := testMapField[*uint](t, fields, "property5")
	property6 := testMapField[*uint64](t, fields, "property6")
	property7 := testMapField[*uint32](t, fields, "property7")
	property8 := testMapField[*uint16](t, fields, "property8")
	property9 := testMapField[*uint8](t, fields, "property9")

	ptr.Property1 = "abc"
	ptr.Property2 = 1
	ptr.Property3 = 2.3
	ptr.Property4 = true
	ptr.Property5 = 4
	ptr.Property6 = 5
	ptr.Property7 = 6
	ptr.Property8 = 7
	ptr.Property9 = 8

	testValueEquals(t, *property1, "abc")
	testValueEquals(t, *property2, 1)
	testValueEquals(t, *property3, 2.3)
	testValueEquals(t, *property4, true)
	testValueEquals(t, *property5, 4)
	testValueEquals(t, *property6, 5)
	testValueEquals(t, *property7, 6)
	testValueEquals(t, *property8, 7)
	testValueEquals(t, *property9, 8)
}

func testMapField[TFieldType interface{}](t *testing.T, fields map[string]any, fieldName string) TFieldType {
	var out TFieldType
	t.Helper()
	field, ok := fields[fieldName]
	if !ok {
		t.Errorf("%s field not found", fieldName)
		t.FailNow()
		return out
	}

	typedField, ok := field.(TFieldType)
	if !ok {
		t.Errorf("%s field type is not correct", fieldName)
		t.FailNow()
		return out
	}

	return typedField
}

func testValueEquals[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
		t.FailNow()
	}
}
