package database

import (
	"database/sql"
	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"
	// _ "modernc.org/sqlite"
)

// Config - Used to provide connection details to [SetupDatabase] function
type Config struct {
	// Hostname - Specifies the hostname for connecting to the database.
	Hostname string
	// Port - Port to user when connecting to the database
	Port string
	// Database - Name of the database to connect to.
	Database string
	// Username - Username to user when connecting to the database.
	Username string
	// Password - Specifies the password to use when connecting to the database.
	Password string
	// Driver - Specifies the driver to use when connecting to the database.
	Driver DriverType
}

type DriverType string

const (
	// MySql driver name for [go-sql-driver/mysql package]
	//
	// [go-sql-driver/mysql package]: https://pkg.go.dev/github.com/go-sql-driver/mysql
	MySql DriverType = "mysql"

	// MariaDB alias driver name for [go-sql-driver/mysql package]
	//
	// [go-sql-driver/mysql package]: https://pkg.go.dev/github.com/go-sql-driver/mysql
	MariaDB DriverType = "mysql"

	// SQLite driver name for [modernc.org/sqlite package]
	//
	// [modernc.org/sqlite package]: https://pkg.go.dev/modernc.org/sqlite
	SQLite DriverType = "sqlite"

	// Postgres driver name for [lib/pq package]
	//
	// [lib/pq package]: https://pkg.go.dev/github.com/lib/pq
	Postgres DriverType = "postgres"
)

// Database contains the [sql.DB] database connection and basic database info.
type Database struct {
	Db           *sql.DB
	dbName       string
	tableColumns map[string][]string
}

// SetupDatabase - Initial setup function for goApiDb connection.
//
// Uses [sql.Open] and info from [Config] parameter to establish database connection.
//
// Returns [Database] connection.
func SetupDatabase(config Config) (*Database, error) {
	db, err := sql.Open(string(config.Driver), getConnectionString(config))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	output := &Database{Db: db, dbName: config.Database, tableColumns: make(map[string][]string)}

	setupTableVariables(output)

	return output, nil
}

func setupTableVariables(database *Database) {
	if database.tableColumns == nil {
		database.tableColumns = make(map[string][]string)
	}
}
