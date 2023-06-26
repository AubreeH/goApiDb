package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Driver   DriverType
}

type DriverType string

const (
	MySql    DriverType = "mysql"
	MariaDB  DriverType = "mysql"
	SQLite   DriverType = "sqlite3"
	Postgres DriverType = "pgx"
)

type Database struct {
	Db           *sql.DB
	dbName       string
	tableColumns map[string][]string
}

func SetupDatabase(config Config) (*Database, error) {
	db, err := sql.Open(string(config.Driver), getConnectionString(config))

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

func setupTableVariables(database *Database) {
	if database.tableColumns == nil {
		database.tableColumns = make(map[string][]string)
	}
}
