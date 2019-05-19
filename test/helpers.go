package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	log "github.com/sirupsen/logrus"
)

const API_KEY = "123456"
const DB_NAME = "./test-data/robolucha-api-test.db"

func PerformRequest(r http.Handler, method, path string, body string, authorization string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", authorization)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func PerformRequestNoAuth(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type MockPublisher struct {
	LastMessage string
	LastChannel string
}

func (mock *MockPublisher) Publish(channel string, message string) {
	mock.LastChannel = channel
	mock.LastMessage = message

	log.WithFields(log.Fields{
		"channel": channel,
		"message": message,
	}).Debug("mock publisher")
}
