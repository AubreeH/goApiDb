package database

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
