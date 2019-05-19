package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/test"
)

const TEST_USERNAME = "foo"

func TestCreateMatch(t *testing.T) {
	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	plan, _ := ioutil.ReadFile("test-data/create-match.json")
	body := string(plan)
	fmt.Println(body)

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/match", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateGameComponent(t *testing.T) {
	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	plan, _ := ioutil.ReadFile("test-data/create-gamecomponent1.json")
	body := string(plan)
	fmt.Println(body)

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/game-component", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var luchador Luchador
	json.Unmarshal(w.Body.Bytes(), &luchador)
	assert.True(t, luchador.ID > 0)
	log.WithFields(log.Fields{
		"luchador.ID": luchador.ID,
	}).Info("First call to create game component")

	// retry to check for duplicate
	w = test.PerformRequest(router, "POST", "/internal/game-component", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var luchador2 Luchador
	json.Unmarshal(w.Body.Bytes(), &luchador2)
	assert.True(t, luchador.ID == luchador2.ID)
	log.WithFields(log.Fields{
		"luchador.ID": luchador2.ID,
	}).Info("Second call to create game component")

	luchadorFromDB := dataSource.findLuchadorByID(luchador2.ID)
	log.WithFields(log.Fields{
		"luchador.configs": luchadorFromDB.Configs,
	}).Info("configs from luchador")

	// all the Mask config items should be present
	for _, color := range maskColors {
		found := false
		for _, config := range luchadorFromDB.Configs {
			if config.Key == color {
				found = true
				break
			}
		}
		assert.True(t, found)
		log.WithFields(log.Fields{
			"color": color,
		}).Info("Color found in luchador config")
	}

	for shape, _ := range maskShapes {
		found := false
		for _, config := range luchadorFromDB.Configs {
			if config.Key == shape {
				found = true
				break
			}
		}
		assert.True(t, found)
		log.WithFields(log.Fields{
			"shape": shape,
		}).Info("Shape found in luchador config")
	}

	getConfigs(t, router, luchador2.ID)
	configsFromDB := getConfigs(t, router, luchador2.ID)
	AssertConfigMatch(t, luchadorFromDB.Configs, configsFromDB)
}

func TestAddScores(t *testing.T) {
	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	luchador1 := dataSource.createLuchador(&Luchador{Name: "foo"})
	luchador2 := dataSource.createLuchador(&Luchador{Name: "bar"})
	luchador3 := dataSource.createLuchador(&Luchador{Name: "dee"})

	matchData := Match{
		Duration:        600000,
		MinParticipants: 1,
		MaxParticipants: 10,
		TimeStart:       time.Now(),
		Participants:    []Luchador{*luchador1, *luchador2, *luchador3},
	}

	match := dataSource.createMatch(&matchData)

	matchID := fmt.Sprintf("%v", match.ID)
	luchador1ID := fmt.Sprintf("%v", luchador1.ID)
	luchador2ID := fmt.Sprintf("%v", luchador2.ID)
	luchador3ID := fmt.Sprintf("%v", luchador3.ID)

	plan, _ := ioutil.ReadFile("test-data/add-match-scores.json")
	body := string(plan)

	body = strings.Replace(body, "{{.matchID}}", matchID, -1)
	body = strings.Replace(body, "{{.luchadorID1}}", luchador1ID, -1)
	body = strings.Replace(body, "{{.luchadorID2}}", luchador2ID, -1)
	body = strings.Replace(body, "{{.luchadorID3}}", luchador3ID, -1)

	log.WithFields(log.Fields{
		"body": body,
	}).Info("TestAddScores")

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/add-match-scores", body, test.API_KEY)
	resultScores := dataSource.getMatchScoresByMatchID(match.ID)
	assert.Equal(t, 3, len(*resultScores))
	assert.Equal(t, http.StatusOK, w.Code)

	// parse request body in object to validate the result
	var scoreList ScoreList
	json.Unmarshal([]byte(body), &scoreList)

	// check if all data was saved correctly
	for _, scoreFromBody := range scoreList.Scores {
		found := false
		for _, scoreFromDB := range *resultScores {
			if scoreFromBody.LuchadorID == scoreFromDB.LuchadorID &&
				scoreFromBody.Kills == scoreFromDB.Kills &&
				scoreFromBody.Deaths == scoreFromDB.Deaths &&
				scoreFromBody.Score == scoreFromDB.Score {

				found = true
				break
			}
		}
		assert.True(t, found)
		log.WithFields(log.Fields{
			"score-from-body": scoreFromBody,
		}).Info("TestAddScores")
	}

}

func getConfigs(t *testing.T, router *gin.Engine, id uint) []Config {

	w := test.PerformRequestNoAuth(router, "GET", fmt.Sprintf("/private/mask-config/%v", id), "")
	assert.Equal(t, http.StatusOK, w.Code)

	var configs []Config
	json.Unmarshal(w.Body.Bytes(), &configs)

	return configs
}

func SessionAllwaysValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := dataSource.createUser(User{Username: TEST_USERNAME})
		c.Set("user", user)
	}
}
