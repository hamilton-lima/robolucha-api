package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

const API_KEY = "123456"
const DB_NAME = "./tests/robolucha-api-test.db"

func performRequest(r http.Handler, method, path string, body string, authorization string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", authorization)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateMatch(t *testing.T) {
	os.Setenv("GORM_DEBUG", "true")
	os.Remove(DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(DB_NAME))
	defer dataSource.db.Close()

	plan, _ := ioutil.ReadFile("tests/create-match.json")
	body := string(plan)
	fmt.Println(body)

	router := createRouter(API_KEY, "true")
	w := performRequest(router, "POST", "/internal/match", body, API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateGameComponent(t *testing.T) {
	os.Setenv("GORM_DEBUG", "true")
	os.Setenv("API_ADD_TEST_USERS", "true")

	os.Remove(DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(DB_NAME))
	defer dataSource.db.Close()
	addTestUsers(dataSource)

	plan, _ := ioutil.ReadFile("tests/create-gamecomponent1.json")
	body := string(plan)
	fmt.Println(body)

	router := createRouter(API_KEY, "true")
	w := performRequest(router, "POST", "/internal/game-component", body, API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var luchador Luchador
	json.Unmarshal(w.Body.Bytes(), &luchador)
	assert.True(t, luchador.ID > 0)
	log.WithFields(log.Fields{
		"luchador.ID": luchador.ID,
	}).Info("First call to create game component")

	// retry to check for duplicate
	w = performRequest(router, "POST", "/internal/game-component", body, API_KEY)
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
	elementsMatch(t, luchadorFromDB.Configs, configsFromDB)
}

type MockPublisher struct {
}

func (redis MockPublisher) Publish(channel string, message string) {
	log.WithFields(log.Fields{
		"channel": channel,
		"message": message,
	}).Info("mock publisher")
}

func TestRenameLuchador(t *testing.T) {
	os.Setenv("GORM_DEBUG", "true")
	os.Setenv("API_ADD_TEST_USERS", "true")

	os.Remove(DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(DB_NAME))
	defer dataSource.db.Close()
	addTestUsers(dataSource)

	publisher = MockPublisher{}

	router := createRouter(API_KEY, "true")

	// we have to login to make name changes
	w := performRequest(router, "POST", "/public/login", `{"email": "foo@bar"}`, "")
	assert.Equal(t, http.StatusOK, w.Code)
	var loginResponse LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	log.WithFields(log.Fields{
		"UUID": loginResponse.UUID,
	}).Info("after login")

	w = performRequest(router, "GET", "/luchador", "", loginResponse.UUID)
	assert.Equal(t, http.StatusOK, w.Code)
	var luchador Luchador
	json.Unmarshal(w.Body.Bytes(), &luchador)
	t.Log(luchador)

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("luchador after login")

	assert.False(t, loginResponse.Error)
	assert.Greater(t, len(loginResponse.UUID), 0)

	// first try to change to a valid name
	luchador.Name = "lucharito"
	plan2, _ := json.Marshal(luchador)
	body2 := string(plan2)
	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("luchador before update")

	w = performRequest(router, "PUT", "/private/luchador", body2, loginResponse.UUID)
	assert.Equal(t, http.StatusOK, w.Code)

	var response UpdateLuchadorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"body":     w.Body.String(),
		"response": response,
	}).Info("after luchador update")

	// assert.False(t, response.luchador.Name != prevName)
	assert.Equal(t, "lucharito", response.luchador.Name)

	// then try a too large name
	// luchador.Name = "123456789012345678901234567890aaaaaa"
}

func TestAddScores(t *testing.T) {
	os.Setenv("GORM_DEBUG", "true")
	os.Remove(DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(DB_NAME))
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

	plan, _ := ioutil.ReadFile("tests/add-match-scores.json")
	body := string(plan)

	body = strings.Replace(body, "{{.matchID}}", matchID, -1)
	body = strings.Replace(body, "{{.luchadorID1}}", luchador1ID, -1)
	body = strings.Replace(body, "{{.luchadorID2}}", luchador2ID, -1)
	body = strings.Replace(body, "{{.luchadorID3}}", luchador3ID, -1)

	log.WithFields(log.Fields{
		"body": body,
	}).Info("TestAddScores")

	router := createRouter(API_KEY, "true")
	w := performRequest(router, "POST", "/internal/add-match-scores", body, API_KEY)
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

func elementsMatch(t *testing.T, a []Config, b []Config) {
	for _, configA := range a {
		found := false
		for _, configB := range b {

			if configA.Key == configB.Key && configA.Value == configB.Value {
				found = true
				break
			}
		}
		assert.True(t, found)
		log.WithFields(log.Fields{
			"config": configA,
		}).Info("match found for config")
	}
}

func getConfigs(t *testing.T, router *gin.Engine, id uint) []Config {

	w := performRequest(router, "POST", "/public/login", `{"email": "foo@bar"}`, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response LoginResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response.Error)
	assert.Greater(t, len(response.UUID), 0)

	log.WithFields(log.Fields{
		"session": response.UUID,
	}).Info("logged in")

	w = performRequest(router, "GET", fmt.Sprintf("/private/mask-config/%v", id), "", response.UUID)
	assert.Equal(t, http.StatusOK, w.Code)

	var configs []Config
	json.Unmarshal(w.Body.Bytes(), &configs)

	return configs
}
