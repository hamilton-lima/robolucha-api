package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"testing"

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

func testStartMatch(t *testing.T, router *gin.Engine, luchadorID uint) model.Match {
	setup.CreateAvailableMatches(ds)
	availableMatches := *ds.FindPublicAvailableMatch()

	log.WithFields(log.Fields{
		"availableMatches": availableMatches,
	}).Info("testStartMatch")

	url := fmt.Sprintf("/private/play/%v", availableMatches[0].ID)
	w := test.PerformRequest(router, "POST", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Match
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"match": availableMatches,
	}).Info("after play")

	return response
}

func testEndMatch(t *testing.T, router *gin.Engine, match model.Match) {
	plan, _ := json.Marshal(match)
	body := string(plan)

	url := "/internal/end-match"
	w := test.PerformRequest(router, "PUT", url, body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)
}

func testGetUser(t *testing.T, router *gin.Engine) model.UserDetails {
	url := "/private/get-user"
	w := test.PerformRequest(router, "GET", url, "", test.API_KEY)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.UserDetails
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

func testAddParticipantToMatch(t *testing.T, matchID uint, luchadorID uint) {
	matchParticipant := model.MatchParticipant{
		LuchadorID: luchadorID,
		MatchID:    matchID,
	}

	plan, _ := json.Marshal(matchParticipant)
	body := string(plan)

	url := "/internal/match-participant"
	w := test.PerformRequest(router, "POST", url, body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUserLevelAfterMatchEnds(t *testing.T) {
	luchador := Setup(t)
	defer ds.DB.Close()

	// Create a gamedefinition with unblockLevel == 16
	testCreateGameDefinition()

	// This is the request from game to start the match and generate the event
	// that an user is joining a match
	match := testStartMatch(t, router, luchador.ID)

	// This call would be done by the runner
	testAddParticipantToMatch(t, match.ID, luchador.ID)

	// End match call by the runner that updates the match
	// AND updates match participants user levels
	testEndMatch(t, router, match)

	userDetails := testGetUser(t, router)
	log.WithFields(log.Fields{
		"userDetails": userDetails,
	}).Info("user details after match ends")

	assert.Equal(t, userDetails.Level.Level, uint(16))
}
