package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

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

func testStartMatch(t *testing.T, router *gin.Engine, luchadorID uint, userName string) model.Match {
	setup.CreateAvailableMatches(ds)
	availableMatches := *ds.FindPublicAvailableMatch()

	log.WithFields(log.Fields{
		"availableMatches": availableMatches,
	}).Info("testStartMatch")

	playRequest := model.PlayRequest{
		AvailableMatchID: availableMatches[0].ID,
		TeamID:           0,
	}

	plan, _ := json.Marshal(playRequest)
	body := string(plan)

	url := "/private/play"
	w := test.PerformRequest(router, "POST", url, body, userName)
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

func testGetUser(t *testing.T, router *gin.Engine, userName string) model.UserDetails {
	url := "/private/get-user"
	w := test.PerformRequest(router, "GET", url, "", userName)

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
	userName := "someOtherPlayer"
	luchador := SetupWithUserName(t, userName)
	defer ds.DB.Close()

	// Create a gamedefinition with unblockLevel == 16
	testCreateGameDefinition()

	// This is the request from game to start the match and generate the event
	// that an user is joining a match
	match := testStartMatch(t, router, luchador.ID, userName)

	// This call would be done by the runner
	testAddParticipantToMatch(t, match.ID, luchador.ID)

	// End match call by the runner that updates the match
	// AND updates match participants user levels
	testEndMatch(t, router, match)

	userDetails := testGetUser(t, router, userName)
	log.WithFields(log.Fields{
		"userDetails": userDetails,
	}).Info("user details after match ends")

	assert.Equal(t, userDetails.Level.Level, uint(16))
}
