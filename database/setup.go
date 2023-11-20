package database

import "database/sql"

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
