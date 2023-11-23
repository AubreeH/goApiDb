package driver

func (d *DriverType) GetConnectionString(username, password, port, hostname, database string) string {
	switch *d {
	case MySql:
		return getMySqlConnectionString(username, password, port, hostname, database)
	case Postgres:
		return getPostgresConnectionString(username, password, port, hostname, database)
	case SQLite:
		return getSQLiteConnectionString(hostname)
	default:
		return ""
	}
}
