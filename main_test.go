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

func performRequest(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", API_KEY)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateMatch(t *testing.T) {
	os.Setenv("GORM_DEBUG", "false")

	dataSource = NewDataSource(BuildSQLLiteConfig("./tests/robolucha-api-test.db"))
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

	dataSource = NewDataSource(BuildSQLLiteConfig("./tests/robolucha-api-test.db"))
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

}
