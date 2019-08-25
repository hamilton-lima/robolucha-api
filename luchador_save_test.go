package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/test"
)

var router *gin.Engine
var mockPublisher *test.MockPublisher

func Setup(t *testing.T) *model.GameComponent {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	os.Setenv("GIN_MODE", "release")

	err := os.Remove(test.DB_NAME)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error removing TEST database")
	}
	ds = datasource.NewDataSource(datasource.BuildSQLLiteConfig(test.DB_NAME))

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher

	router = createRouter(test.API_KEY, "true", SessionAllwaysValid)

	luchador := GetLuchador(t)
	return &luchador
}

func GetLuchador(t *testing.T) model.GameComponent {
	getLuchador := test.PerformRequestNoAuth(router, "GET", "/private/luchador", "")
	var luchador model.GameComponent
	json.Unmarshal(getLuchador.Body.Bytes(), &luchador)
	return luchador
}

func TestLuchadorUpdateDuplicatedNameSameUser(t *testing.T) {
	luchador := Setup(t)
	defer ds.DB.Close()

	plan2, _ := json.Marshal(luchador)
	body2 := string(plan2)

	log.WithFields(log.Fields{
		"luchador": luchador.Name,
	}).Debug("luchador before same name update")

	w := test.PerformRequestNoAuth(router, "PUT", "/private/luchador", body2)
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.UpdateLuchadorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response":        response.Luchador.Name,
		"response.errors": response.Errors,
	}).Debug("after luchador update")

	t.Log(response.Errors)
	assert.Equal(t, len(response.Errors), 0)

}

func TestLuchadorUpdateLongName(t *testing.T) {
	luchador := Setup(t)
	defer ds.DB.Close()

	luchador.Name = "123456789 123456789 123456789 123456789 A"
	response := UpdateLuchador(t, router, luchador)
	assert.True(t, len(response.Errors) > 0)
}

func TestLuchadorUpdateEmptyAndSmallNames(t *testing.T) {
	luchador := Setup(t)
	defer ds.DB.Close()

	// then try a too large name
	luchador.Name = "A"
	response := UpdateLuchador(t, router, luchador)
	assert.True(t, len(response.Errors) > 0)
}

func TestLuchadorUpdateName(t *testing.T) {
	luchador := Setup(t)
	defer ds.DB.Close()

	// first try to change to a valid name
	luchador.Name = "lucharito"
	response := UpdateLuchador(t, router, luchador)
	assert.Equal(t, "lucharito", response.Luchador.Name)
	assert.Equal(t, 0, len(response.Errors))

	channel := fmt.Sprintf("luchador.%v.update", luchador.ID)

	log.WithFields(log.Fields{
		"expected":         channel,
		"publishedChannel": mockPublisher.LastChannel,
	}).Debug("publish event")
	assert.True(t, mockPublisher.LastChannel == channel)

}
func TestLuchadorUpdateRandomMask(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	luchador := Setup(t)
	defer ds.DB.Close()

	// assign new random Configs to update the luchador
	var originalConfigs []model.Config = luchador.Configs
	updatedConfigs := make([]model.Config, len(originalConfigs))

	for n, config := range originalConfigs {
		updatedConfigs[n].Key = config.Key
		updatedConfigs[n].Value = config.Value + "A"
	}

	luchador.Configs = updatedConfigs

	plan2, _ := json.Marshal(luchador)
	body2 := string(plan2)

	w := test.PerformRequestNoAuth(router, "PUT", "/private/luchador", body2)
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.UpdateLuchadorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response.Errors":   response.Errors,
		"response.Luchador": response.Luchador,
	}).Debug("after luchador update")

	// check if no errors exist in the response
	assert.Equal(t, 0, len(response.Errors))

	// check if configs are updated in the response
	assert.Equal(t, len(updatedConfigs), len(response.Luchador.Configs))
	AssertConfigMatch(t, updatedConfigs, response.Luchador.Configs)

	// check if configs are updated in the subsequent GET of luchador
	afterUpdateLuchador := GetLuchador(t)
	assert.Equal(t, len(updatedConfigs), len(afterUpdateLuchador.Configs))
	AssertConfigMatch(t, updatedConfigs, afterUpdateLuchador.Configs)

	// check if after update the correct event is published
	channel := fmt.Sprintf("luchador.%v.update", luchador.ID)

	log.WithFields(log.Fields{
		"expected":         channel,
		"publishedChannel": mockPublisher.LastChannel,
	}).Debug("publish event")
	assert.True(t, mockPublisher.LastChannel == channel)

}
