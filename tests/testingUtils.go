package tests

import (
	"fmt"
	"github.com/AubreeH/goApiDb/database"
	"github.com/AubreeH/goApiDb/helpers"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"strings"
	"time"
)

var db *database.Database

func init() {
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load("../.env")
	conf := getDatabaseConfig()
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

	return closeFunc, database.BuildTable(db, entity)
}

func seedTable(count int, table string, columns map[string]string) error {
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
	for i := 0; i < count; i++ {
		values = append(values, "("+strings.Join(valuePlaceholders, ", ")+")")

		for _, v := range types {
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
		}
	}

	valuesString := strings.Join(values, ", ")
	columnsString := strings.Join(columnNames, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, columnsString, valuesString)
	_, err := db.Db.Exec(query, args...)
	return err
}

func seedTableWithValueInMiddle(count int, table string, columns map[string]string, data map[string]any) (int64, error) {
	lowerBound := rand.Intn(count-19) + 10
	upperBound := count - lowerBound - 1

	err := seedTable(lowerBound, table, columns)
	if err != nil {
		return 0, err
	}

	id, err := createDatabaseRow(db, table, data)
	if err != nil {
		return 0, err
	}

	err = seedTable(upperBound, table, columns)
	if err != nil {
		return 0, err
	}

	return id, nil

}
