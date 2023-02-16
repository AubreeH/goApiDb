package database

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/helpers"
	"log"
	"os"
	"reflect"
	"time"
)

var db *sql.DB
var tableColumns map[string][]string

func SetupDatabase() *sql.DB {
	if db != nil {
		return db
	}

	for true {
		connectionString := getConnectionString()

		var err error
		db, err = sql.Open("mysql", connectionString)

		if err != nil {
			log.Print("Error whilst opening connection to db", err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Print("Error whilst pinging connection to db", err)
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	SetupTableVariables()

	return db
}

func GetDb() *sql.DB {
	return db
}

func getConnectionString() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	var account string
	if pass != "" {
		account = user + ":" + pass
	} else {
		account = user
	}

	url := host + ":" + port

	return account + "@tcp(" + url + ")/" + name + "?parseTime=true"
}

func SetupTableVariables() {
	if tableColumns == nil {
		tableColumns = make(map[string][]string)
	}
}

func getEntityConstruction[T any](entity *T) (map[string]any, T, string) {
	val := reflect.ValueOf(entity).Elem()

	tmp := reflect.New(val.Elem().Type()).Elem()
	tmp.Set(val.Elem())

	columnVariables := make(map[string]any)
	getColumnsFromStruct(tmp, columnVariables)
	tableName := helpers.GetTableName(reflect.ValueOf(entity).Elem().Interface())

	return columnVariables, tmp.Addr().Interface().(T), tableName
}

func getColumnsFromStruct(refValue reflect.Value, columnVariables map[string]any) map[string]any {

	numFields := refValue.NumField()
	for i := 0; i < numFields; i++ {
		if helpers.ParseBool(refValue.Type().Field(i).Tag.Get("sql_ignore")) {
			continue
		}

		valueField := refValue.Field(i)
		getPtrFunc := valueField.MethodByName("GetPtrFunc")

		sqlName := refValue.Type().Field(i).Tag.Get("sql_name")
		if getPtrFunc.IsValid() {
			columnVariables[sqlName] = getPtrFunc.Call([]reflect.Value{valueField.Addr()})[0].Interface()
		} else if valueField.Kind().String() == "struct" && refValue.Type().Field(i).Tag.Get("parse_struct") != "false" {
			getColumnsFromStruct(valueField, columnVariables)
		} else {
			columnVariables[sqlName] = valueField.Addr().Interface()
		}
	}
	return columnVariables
}

func BuildRow(entity interface{}, result *sql.Rows) ([]interface{}, interface{}, error) {
	columnVariables, ptr, tableName := getEntityConstruction(&entity)

	var columns []string
	if tableColumns[tableName] != nil {
		columns = tableColumns[tableName]
	} else {
		resultColumns, err := result.Columns()
		if err != nil {
			return nil, ptr, err
		}
		tableColumns[tableName] = resultColumns
		columns = resultColumns
	}

	retArgs := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		retArgs[i] = columnVariables[columns[i]]
	}

	return retArgs, ptr, nil
}
