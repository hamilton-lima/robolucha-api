package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/test"

	log "github.com/sirupsen/logrus"
)

func UpdateLuchador(t *testing.T, router *gin.Engine, session string, luchador *Luchador) UpdateLuchadorResponse {
	plan, _ := json.Marshal(luchador)
	body := string(plan)

	log.WithFields(log.Fields{
		"luchador": luchador.Name,
	}).Info("luchador before update")

	w := test.PerformRequest(router, "PUT", "/private/luchador", body, session)
	assert.Equal(t, http.StatusOK, w.Code)

	var response UpdateLuchadorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response": response,
	}).Info("after luchador update")

	return response
}

func AssertConfigMatch(t *testing.T, a []Config, b []Config) {
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
			"found":  found,
		}).Info("match found for config")
	}
}

func CountChangesConfigMatch(t *testing.T, a []Config, b []Config) int {
	counter := 0
	for _, configA := range a {
		found := false
		for _, configB := range b {

			if configA.Key == configB.Key && configA.Value != configB.Value {
				found = true
				break
			}
		}

		if !found {
			counter++
		}
	}

	return counter
}
