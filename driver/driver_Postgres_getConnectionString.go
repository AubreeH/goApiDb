package driver

import "fmt"

func getPostgresConnectionString(username, password, port, hostname, database string) string {
	var account string
	if password != "" {
		account = username + ":" + password
	} else {
		account = username
	}

	var url string
	if port != "" {
		url = hostname + ":" + port
	} else {
		url = hostname
	}

	return fmt.Sprintf("postgresql://%s@%s/%s?sslmode=disable", account, url, database)
}
