package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	log "github.com/sirupsen/logrus"
)

const API_KEY = "123456"
const DB_NAME = "robolucha-api-test.db"

// PerformRequest sends http request no authentication header
func PerformRequest(r http.Handler, method, path string, body string, authorization string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", authorization)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// PerformRequestNoAuth sends http request no authentication
func PerformRequestNoAuth(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// MockPublisher definition
type MockPublisher struct {
	LastMessage string
	LastChannel string
	Messages    map[string][]string
}

// Publish saves the messages in the memory and in lastMessage/Channel
func (mock *MockPublisher) Publish(channel string, message string) {
	mock.LastChannel = channel
	mock.LastMessage = message

	log.WithFields(log.Fields{
		"channel": channel,
		"message": message,
	}).Info("mock publisher")

	if mock.Messages == nil {
		mock.Messages = make(map[string][]string)
	}

	mock.Messages[channel] = append(mock.Messages[channel], message)
}

// ResetMessages clear previous messages
func (mock *MockPublisher) ResetMessages() {
	mock.Messages = make(map[string][]string)
}
