package tests

import (
	"fmt"
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/entities"
	"github.com/AubreeH/goApiDb/helpers"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var db *database.Database

func init() {
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
}

func InitDb() {
	conf := getDatabaseConfig()
	var err error
	db, err = database.SetupDatabase(conf)
	if err != nil {
		panic(err)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getDatabaseConfig() database.DatabaseConfig {
	return database.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func createDatabaseRow(db *database.Database, table string, data map[string]any) (int64, error) {
	var columns []string
	var values []any
	var valuePlaceholders []string
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
		valuePlaceholders = append(valuePlaceholders, "?")
	}

	columnsStr := strings.Join(columns, ", ")
	valuesStr := strings.Join(valuePlaceholders, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columnsStr, valuesStr)
	result, err := db.Db.Exec(query, values...)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertId, err
}

func setupTable(entity interface{}) (func(), error) {
	closeFunc := func() {
		tableName := helpers.GetTableName(entity)
		_, err := db.Db.Exec("DROP TABLE " + tableName)
		if err != nil {
			panic(err)
		}
		return
	}

	return closeFunc, database.BuildTable(db, entity, false, true)
}

func dropTable[T any]() error {
	var entity T
	tableInfo, err := entities.GetTableInfo(entity)
	if err != nil {
		return err
	}
	_, err = db.Db.Exec("DROP TABLE " + tableInfo.Name)
	if err != nil {
		return err
	}
	return nil
}

func seedTable(count int, table string, columns map[string]string) (map[int64]map[string]any, error) {
	var columnNames []string
	var valuePlaceholders []string
	var types []string

	for k, v := range columns {
		columnNames = append(columnNames, k)
		valuePlaceholders = append(valuePlaceholders, "?")
		types = append(types, v)
	}

	var values []string
	var args []any

	var seededValues = make(map[int64]map[string]any, count+1)
	for i := 0; i < count; i++ {
		values = append(values, "("+strings.Join(valuePlaceholders, ", ")+")")

		seededValues[int64(i+1)] = make(map[string]any, len(columnNames))

		for k, v := range types {
			var arg any
			switch v {
			case "string":
				arg = randSeq(20)
				break
			case "number":
				arg = rand.Int()
				break
			case "date":
				arg = time.UnixMilli(rand.Int63n(time.Now().UnixMilli()))
			}
			args = append(args, arg)
			seededValues[int64(i+1)][columnNames[k]] = arg
		}
	}

	valuesString := strings.Join(values, ", ")
	columnsString := strings.Join(columnNames, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, columnsString, valuesString)
	_, err := db.Db.Exec(query, args...)
	return seededValues, err
}

func seedTableWithValueInMiddle(count int, table string, columns map[string]string, data map[string]any) (int64, map[int64]map[string]any, error) {
	lowerBound := rand.Intn(count-19) + 10
	upperBound := count - lowerBound - 1

	var output map[int64]map[string]any

	seededValuesLower, err := seedTable(lowerBound, table, columns)
	if err != nil {
		return 0, output, err
	}

	id, err := createDatabaseRow(db, table, data)
	if err != nil {
		return 0, output, err
	}

	seededValuesUpper, err := seedTable(upperBound, table, columns)
	if err != nil {
		return 0, output, err
	}

	for k, v := range seededValuesUpper {
		seededValuesLower[k] = v
	}

	return id, seededValuesLower, nil

}

type c struct {
	Condition bool
	Args      []any
}

func condition(condition bool, args ...any) c {
	return c{
		Condition: condition,
		Args:      args,
	}
}

func assert(t *testing.T, conditions ...c) {
	fail := false

	for _, v := range conditions {
		if v.Condition {
			fail = true
			t.Error(v.Args...)
		}
	}

	if fail {
		t.FailNow()
	}
}

func assertError(t *testing.T, err error) {
	_, file, line, _ := runtime.Caller(1)
	assert(t, condition(err != nil, fmt.Sprintf("Error in %s on line %d: ", file, line), err))

}

func p(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	output := []any{fmt.Sprintf("%s:%d: ", filepath.Base(file), line)}
	output = append(output, args...)
	fmt.Print(output...)
}
