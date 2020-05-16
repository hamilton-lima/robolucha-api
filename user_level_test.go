package main

import (
	"encoding/json"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/test"
)

func GetUserDetails(t *testing.T) model.UserDetails {
	jsonData := test.PerformRequestNoAuth(router, "GET", "/private/get-user", "")
	var userDetails model.UserDetails
	json.Unmarshal(jsonData.Body.Bytes(), &userDetails)
	return userDetails
}

func TestUserLevelDefault(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	firstDetails := GetUserDetails(t)
	log.WithFields(log.Fields{
		"userDetails": firstDetails,
	}).Debug("First call to get user details")

	assert.Equal(t, uint(0), firstDetails.Level.Level)
}

func TestUserLevelAfterUpdate(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	firstDetails := GetUserDetails(t)
	log.WithFields(log.Fields{
		"userDetails": firstDetails,
	}).Debug("First call to get user details")

	assert.Equal(t, uint(0), firstDetails.Level.Level)

	firstDetails.Level.Level = 42
	ds.UpdateUserLevel(&firstDetails.Level)

	secondDetails := GetUserDetails(t)
	log.WithFields(log.Fields{
		"userDetails": secondDetails,
	}).Debug("***** Second call to get user details")

	assert.Equal(t, uint(42), secondDetails.Level.Level)

}
