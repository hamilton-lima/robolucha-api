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

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

const API_KEY = "123456"
const DB_NAME = "./tests/robolucha-api-test.db"

func performRequest(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", API_KEY)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateMatch(t *testing.T) {
	os.Setenv("GORM_DEBUG", "false")
	os.Remove(DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(DB_NAME))
	defer dataSource.db.Close()

	plan, _ := ioutil.ReadFile("tests/create-match.json")
	body := string(plan)
	fmt.Println(body)

	router := createRouter(API_KEY, "true")
	w := performRequest(router, "POST", "/internal/match", body)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateGameComponent(t *testing.T) {
	os.Setenv("GORM_DEBUG", "false")
	os.Remove(DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(DB_NAME))
	defer dataSource.db.Close()

	plan, _ := ioutil.ReadFile("tests/create-gamecomponent1.json")
	body := string(plan)
	fmt.Println(body)

	router := createRouter(API_KEY, "true")
	w := performRequest(router, "POST", "/internal/game-component", body)
	assert.Equal(t, http.StatusOK, w.Code)

	var luchador Luchador
	json.Unmarshal(w.Body.Bytes(), &luchador)
	assert.True(t, luchador.ID > 0)
	log.WithFields(log.Fields{
		"luchador.ID": luchador.ID,
	}).Info("First call to create game component")

	// retry to check for duplicate
	w = performRequest(router, "POST", "/internal/game-component", body)
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
	assert.True(t, len(luchadorFromDB.Configs) == len(MASK_CONFIG_KEYS), "Same amount of keys in the config")

	// all the Mask config items should be present
	for _, key := range MASK_CONFIG_KEYS {
		found := false
		for _, config := range luchadorFromDB.Configs {
			if config.Key == key {
				found = true
				break
			}
		}
		assert.True(t, found)
		log.WithFields(log.Fields{
			"key": key,
		}).Info("Key found in luchador config")

	}

}
