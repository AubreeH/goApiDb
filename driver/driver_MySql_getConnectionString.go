package driver

import "fmt"

func getMySqlConnectionString(username, password, port, hostname, database string) string {
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

	return fmt.Sprintf("%s@%s(%s)/%s?parseTime=true", account, "tcp", url, database)
}
