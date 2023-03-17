package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Database struct {
	Db           *sql.DB
	dbName       string
	tableColumns map[string][]string
}

func SetupDatabase(config Config) (*Database, error) {
	connectionString := getConnectionString(config)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	output := &Database{Db: db, dbName: config.Name, tableColumns: make(map[string][]string)}

	setupTableVariables(output)

	return output, nil
}

func getConnectionString(config Config) string {

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

func setupTableVariables(database *Database) {
	if database.tableColumns == nil {
		database.tableColumns = make(map[string][]string)
	}
}
