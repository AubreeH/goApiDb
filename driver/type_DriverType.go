package driver

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
