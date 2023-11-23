package database

import (
	"database/sql"

	"github.com/AubreeH/goApiDb/driver"
)

// SetupDatabase - Initial setup function for goApiDb connection.
//
// Uses [sql.Open] and info from [Config] parameter to establish database connection.
//
// Returns [Database] connection.
func SetupDatabase(config Config) (*Database, error) {
	db, err := sql.Open(string(config.Driver), config.GetConnectionString())

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	output := &Database{Db: db, dbName: config.Database, tableColumns: make(map[string][]string), config: config}

	setupTableVariables(output)

	return output, nil
}

func setupTableVariables(database *Database) {
	if database.tableColumns == nil {
		database.tableColumns = make(map[string][]string)
	}
}

func (d *Database) DescribeTable(tableName string) (*driver.TableDescription, error) {
	return d.config.Driver.DescribeTable(d.Db, tableName)
}

func (d *Database) GetKey(tableName, columnName string) (*driver.Key, error) {
	return d.config.Driver.GetKey(d.Db, tableName, columnName)
}
