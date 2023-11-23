package database

func BuildTables(db *Database, entities ...interface{}) error {
	tableQueries, addConstraintsQueries, dropConstraintsQueries, err := GetUpdateTableQueriesForEntities(db, entities...)
	if err != nil {
		return err
	}

	for _, query := range dropConstraintsQueries {
		if query != "" {
			_, err = db.Db.Exec(query)
			if err != nil {
				return err
			}
		}
	}

	for _, query := range tableQueries {
		if query != "" {
			_, err = db.Db.Exec(query)
			if err != nil {
				return err
			}
		}
	}

	for _, query := range addConstraintsQueries {
		if query != "" {
			_, err = db.Db.Exec(query)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
