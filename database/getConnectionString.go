package database

import (
	"fmt"
)

func getConnectionString(config Config) string {
	switch config.Driver {
	case MySql:
		return getMySqlConnectionString(config)
	case SQLite:
		return getSQLiteConnectionString(config)
	case Postgres:
		return getPostgresConnectionString(config)
	}

	return ""
}

func getMySqlConnectionString(config Config) string {
	var account string
	if config.Password != "" {
		account = config.Username + ":" + config.Password
	} else {
		account = config.Username
	}

	var url string
	if config.Port != "" {
		url = config.Hostname + ":" + config.Port
	} else {
		url = config.Hostname
	}

	return fmt.Sprintf("%s@%s(%s)/%s?parseTime=true", account, "tcp", url, config.Database)
}

func getPostgresConnectionString(config Config) string {
	var account string
	if config.Password != "" {
		account = config.Username + ":" + config.Password
	} else {
		account = config.Username
	}

	var url string
	if config.Port != "" {
		url = config.Hostname + ":" + config.Port
	} else {
		url = config.Hostname
	}

	return fmt.Sprintf("postgresql://%s@%s/%s?sslmode=disable", account, url, config.Database)
}

func getSQLiteConnectionString(config Config) string {
	return config.Hostname
}
