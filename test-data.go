package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func addTestUsers(dataSource *DataSource) {
	apiAddTestUsers := os.Getenv("API_ADD_TEST_USERS")

	log.WithFields(log.Fields{
		"API_ADD_TEST_USERS": apiAddTestUsers,
	}).Debug("addTestUsers")

	if apiAddTestUsers == "true" {
		dataSource.createUser(&User{Email: "foo@bar", Password: "foobar"})
	}
}
