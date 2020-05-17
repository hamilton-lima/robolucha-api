package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gitlab.com/robolucha/robolucha-api/auth"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/setup"
	"gitlab.com/robolucha/robolucha-api/test"
	"gotest.tools/assert"
)

func testCreateGameDefinition() {
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "FOOBAR"
	gd.UnblockLevel = 16
	ds.CreateGameDefinition(&gd)
}

func testStartMatch(t *testing.T, luchadorID uint) model.Match {
	setup.CreateAvailableMatches(ds)
	availableMatches := *ds.FindPublicAvailableMatch()

	url := fmt.Sprintf("/private/play/%v", availableMatches[0].ID)
	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Match
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

func testEndMatch(t *testing.T, match model.Match) {
	plan, _ := json.Marshal(match)
	body := string(plan)

	url := "/internal/end-match"
	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequest(router, "PUT", url, body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)
}

func testGetUser(t *testing.T) model.UserDetails {
	url := "/private/get-user"
	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequest(router, "GET", url, "", test.API_KEY)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.UserDetails
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

func TestUpdateUserLevelAfterMatchEnds(t *testing.T) {
	SetupMain(t)
	defer ds.DB.Close()

	luchador := Setup(t)

	testCreateGameDefinition()
	match := testStartMatch(t, luchador.ID)

	testEndMatch(t, match)
	userDetails := testGetUser(t)
	assert.Equal(t, userDetails.Level.Level, 16)
}
