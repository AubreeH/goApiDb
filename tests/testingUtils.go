package tests

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/driver"
	"github.com/AubreeH/goApiDb/structParsing"
	"github.com/joho/godotenv"
)

var db *database.Database

func init() {
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
}

func InitDb(t *testing.T) {
	t.Helper()
	conf := getDatabaseConfig()
	var err error
	db, err = database.SetupDatabase(conf)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Db.Ping()
	if err != nil {
		t.Fatal(err)
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

func getDatabaseConfig() database.Config {
	return database.Config{
		Hostname: os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Driver:   driver.DriverType(os.Getenv("DB_DRIVER")),
	}
}

func setupTables(t *testing.T, doCleanup bool, ent ...interface{}) {
	t.Helper()
	t.Log("Setting up tables - START")
	assertError(t, database.BuildTables(db, ent...))
	if doCleanup {
		t.Cleanup(func() {
			cleanupTables(t, ent...)
		})
	}
	t.Log("Setting up tables - FINISH")
}

func cleanupTables(t *testing.T, ent ...interface{}) {
	t.Helper()
	t.Log("Cleaning up tables - START")
	for _, e := range ent {
		tableInfo, err := structParsing.GetTableInfo(e)
		assertError(t, err)
		_, err = db.Db.Exec(fmt.Sprintf("DROP TABLE `%s`", tableInfo.Name))
		assertError(t, err)
	}
	t.Log("Cleaning up tables - FINISH")
}

func dropTable[T any]() error {
	var entity T
	tableInfo, err := structParsing.GetTableInfo(entity)
	if err != nil {
		return err
	}
	_, err = db.Db.Exec(fmt.Sprintf("DROP TABLE `%s`", tableInfo.Name))
	if err != nil {
		return err
	}
	return nil
}

func seedTable(t *testing.T, count int, table string, columns map[string]string, data ...map[string]any) ([]int, map[int]map[string]any) {
	t.Helper()
	t.Logf("Seeding table - %s - START", table)
	if count < len(data) {
		assertError(t, fmt.Errorf("count must be greater than or equal to the number of data sets provided"))
	}

	tempMap := make(map[int]map[string]any)
	providedDataIds := make([]int, len(data))

	idColumn := ""
	for columnName, columnType := range columns {
		if columnType == "id" {
			idColumn = columnName
			break
		}
	}

	for i, v := range data {
		id := 0
		for id == 0 || tempMap[id] != nil {
			id = rand.Intn(count) + 1
		}
		providedDataIds[i] = id
		if idColumn != "" {
			v[idColumn] = id
			data[i] = v
		}
		tempMap[id] = v
	}

	seededValues := createSeededEntries(count, columns)
	for i, v := range tempMap {
		seededValues[i] = v
	}

	saveSeededEntriesToDb(t, table, seededValues)

	t.Logf("Seeding table - %s - FINISH", table)
	return providedDataIds, seededValues

}

func createSeededEntries(count int, columns map[string]string) map[int]map[string]any {
	values := make(map[int]map[string]any)
	for i := 0; i < count; i++ {
		values[i+1] = make(map[string]any)
		for columnName, columnType := range columns {
			values[i+1][columnName] = getValueOfType(i+1, columnType)
		}
	}
	return values
}

func saveSeededEntriesToDb(t *testing.T, table string, values map[int]map[string]any) {
	t.Helper()
	t.Logf("Saving seeded entries to db - %s - START", table)

	var columnNames []string
	var valuePlaceholders []string
	var args []any

	for _, v := range values {
		for columnName := range v {
			columnNames = append(columnNames, columnName)
			valuePlaceholders = append(valuePlaceholders, "?")
		}
		break
	}

	valuePlaceholdersString := strings.Join(valuePlaceholders, ", ")
	var valuesStringArr []string

	for _, v := range values {
		for _, columnName := range columnNames {
			args = append(args, v[columnName])
		}
		valuesStringArr = append(valuesStringArr, "("+valuePlaceholdersString+")")
	}

	query := fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES %s", table, strings.Join(columnNames, "`, `"), strings.Join(valuesStringArr, ",\n"))
	_, err := db.Db.Exec(query, args...)
	assertError(t, err)
	t.Logf("Saving seeded entries to db - %s - FINISH", table)
}

func getValueOfType(id any, t string) any {
	switch t {
	case "string":
		return randSeq(20)
	case "number":
		return rand.Int()
	case "date":
		return time.UnixMilli(rand.Int63n(time.Now().UnixMilli())).Format("2006-01-02")
	case "time":
		return time.UnixMilli(rand.Int63n(time.Now().UnixMilli())).Format("15:04:05")
	case "datetime":
		return time.UnixMilli(rand.Int63n(time.Now().UnixMilli())).Format("2006-01-02 15:04:05")
	case "boolean":
		return rand.Intn(2) == 1
	case "id":
		return id
	case "null":
		return nil
	}
	return nil
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
			_, file, line, _ := runtime.Caller(1)
			output := []any{fmt.Sprintf("%s:%d: ", filepath.Base(file), line)}
			output = append(output, v.Args...)
			t.Error(output...)
		}
	}

	if fail {
		t.FailNow()
	}
}

func assertError(t *testing.T, err error) {
	assert(t, e(err))
}

func e(err error) c {
	_, file, line, _ := runtime.Caller(1)
	return condition(err != nil, fmt.Sprintf("Error in %s on line %d: ", file, line), err)
}
