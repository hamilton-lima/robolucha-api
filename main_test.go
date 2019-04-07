package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const API_KEY = "123456"

func setup() {

}
func performRequest(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {

	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", API_KEY)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateMatch(t *testing.T) {
	os.Setenv("GORM_DEBUG", "true")

	dataSource = NewDataSource(BuildSQLLiteConfig("./tests/robolucha-api-test.db"))
	defer dataSource.db.Close()

	plan, _ := ioutil.ReadFile("tests/default-gamedefinition.json")
	json := string(plan)
	fmt.Println(json)

	router := createRouter(API_KEY)

	w := performRequest(router, "POST", "/internal/match", json)
	assert.Equal(t, http.StatusOK, w.Code)
}
