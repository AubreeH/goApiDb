package database

import (
	"database/sql"
	"github.com/AubreeH/goApiDb/helpers"
	"reflect"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Database struct {
	Db           *sql.DB
	tableColumns map[string][]string
}

func SetupDatabase(config DatabaseConfig) (*Database, error) {
	connectionString := getConnectionString(config)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	output := &Database{Db: db}

	SetupTableVariables(output)

	return output, nil
}

func getConnectionString(config DatabaseConfig) string {

	var account string
	if config.Password != "" {
		account = config.User + ":" + config.Password
	} else {
		account = config.User
	}

	var url string
	if config.Port != "" {
		url = config.Host + ":" + config.Port
	} else {
		url = config.Host
	}

	return account + "@tcp(" + url + ")/" + config.Name + "?parseTime=true"
}

func SetupTableVariables(database *Database) {
	if database.tableColumns == nil {
		database.tableColumns = make(map[string][]string)
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

func BuildRow(db *Database, entity interface{}, result *sql.Rows) ([]interface{}, interface{}, error) {
	columnVariables, ptr, tableName := getEntityConstruction(&entity)

	var columns []string
	if db.tableColumns[tableName] != nil {
		columns = db.tableColumns[tableName]
	} else {
		resultColumns, err := result.Columns()
		if err != nil {
			return nil, ptr, err
		}
		db.tableColumns[tableName] = resultColumns
		columns = resultColumns
	}

	retArgs := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		retArgs[i] = columnVariables[columns[i]]
	}

	return retArgs, ptr, nil
}