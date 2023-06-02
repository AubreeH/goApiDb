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

	return fmt.Sprintf("%s@%s(%s)/%s?parseTime=true", account, "tcp", url, config.Name)
}

func getPostgresConnectionString(config Config) string {
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

	return fmt.Sprintf("postgresql://%s@%s/%s?sslmode=disable", account, url, config.Name)
}

func getSQLiteConnectionString(config Config) string {
	return fmt.Sprintf("file:%s?mode=memory&cache=%s", config.Host, "shared")
}
