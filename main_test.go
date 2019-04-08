package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
