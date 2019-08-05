package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/test"

	log "github.com/sirupsen/logrus"
)

// UpdateLuchador definition
func UpdateLuchador(t *testing.T, router *gin.Engine, luchador *model.GameComponent) model.UpdateLuchadorResponse {
	plan, _ := json.Marshal(luchador)
	body := string(plan)

	log.WithFields(log.Fields{
		"luchador": luchador.Name,
	}).Info("luchador before update")

	w := test.PerformRequestNoAuth(router, "PUT", "/private/luchador", body)
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.UpdateLuchadorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response": response,
	}).Info("after luchador update")

	return response
}

// AssertConfigMatch definition
func AssertConfigMatch(t *testing.T, a []model.Config, b []model.Config) {
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

// CountChangesConfigMatch definition
func CountChangesConfigMatch(t *testing.T, a []model.Config, b []model.Config) int {
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
