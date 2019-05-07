package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/test"
)

var router *gin.Engine
var mockPublisher *test.MockPublisher
var session string

func Login(router *gin.Engine) string {
	// we have to login to make name changes
	w := test.PerformRequest(router, "POST", "/public/login", `{"email": "foo@bar"}`, "")
	var loginResponse LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	log.WithFields(log.Fields{
		"UUID": loginResponse.UUID,
	}).Info("after login")

	return loginResponse.UUID
}

func Setup(t *testing.T) *Luchador {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	os.Setenv("GORM_DEBUG", "false")
	os.Setenv("API_ADD_TEST_USERS", "true")

	err := os.Remove(test.DB_NAME)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("error removing TEST database")
	}
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	AddTestUsers(dataSource)

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher

	router = createRouter(test.API_KEY, "true")
	session = Login(router)

	luchador := GetLuchador(t, session)
	return &luchador
}

func GetLuchador(t *testing.T, session string) Luchador {
	getLuchador := test.PerformRequest(router, "GET", "/private/luchador", "", session)
	var luchador Luchador
	json.Unmarshal(getLuchador.Body.Bytes(), &luchador)
	return luchador
}

func TestLuchadorUpdateDuplicatedNameSameUser(t *testing.T) {
	luchador := Setup(t)
	defer dataSource.db.Close()

	plan2, _ := json.Marshal(luchador)
	body2 := string(plan2)

	log.WithFields(log.Fields{
		"luchador": luchador.Name,
	}).Info("luchador before same name update")

	w := test.PerformRequest(router, "PUT", "/private/luchador", body2, session)
	assert.Equal(t, http.StatusOK, w.Code)

	var response UpdateLuchadorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response":        response.Luchador.Name,
		"response.errors": response.Errors,
	}).Info("after luchador update")

	t.Log(response.Errors)
	assert.Equal(t, len(response.Errors), 0)

}

func TestLuchadorUpdateLongName(t *testing.T) {
	luchador := Setup(t)
	defer dataSource.db.Close()

	luchador.Name = "123456789 123456789 123456789 123456789 A"
	response := UpdateLuchador(t, luchador)
	assert.Greater(t, len(response.Errors), 0)
}

func TestLuchadorUpdateEmptyAndSmallNames(t *testing.T) {
	luchador := Setup(t)
	defer dataSource.db.Close()

	// then try a too large name
	luchador.Name = "A"
	response := UpdateLuchador(t, luchador)
	assert.Greater(t, len(response.Errors), 0)
}

func TestLuchadorUpdateName(t *testing.T) {
	luchador := Setup(t)
	defer dataSource.db.Close()

	// first try to change to a valid name
	luchador.Name = "lucharito"
	response := UpdateLuchador(t, luchador)
	assert.Equal(t, "lucharito", response.Luchador.Name)
	assert.Equal(t, 0, len(response.Errors))

	channel := fmt.Sprintf("luchador.%v.update", luchador.ID)

	log.WithFields(log.Fields{
		"expected":         channel,
		"publishedChannel": mockPublisher.LastChannel,
	}).Info("publish event")
	assert.True(t, mockPublisher.LastChannel == channel)

}
func TestLuchadorUpdateRandomMask(t *testing.T) {
	luchador := Setup(t)
	defer dataSource.db.Close()

	// assign new random Configs to update the luchador
	var originalConfigs []Config = luchador.Configs
	var randomConfigs []Config = randomConfig()
	luchador.Configs = randomConfigs

	plan2, _ := json.Marshal(luchador)
	body2 := string(plan2)
	log.WithFields(log.Fields{
		"luchador": luchador.Name,
	}).Info("luchador before update")

	w := test.PerformRequest(router, "PUT", "/private/luchador", body2, session)
	assert.Equal(t, http.StatusOK, w.Code)

	var response UpdateLuchadorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response.Errors":   response.Errors,
		"response.Luchador": response.Luchador,
	}).Info("after luchador update")

	// check if no errors exist in the response
	assert.Equal(t, 0, len(response.Errors))

	// check if configs are updated in the response
	assert.Equal(t, len(randomConfigs), len(response.Luchador.Configs))
	AssertConfigMatch(t, randomConfigs, response.Luchador.Configs)
	changed := CountChangesConfigMatch(t, originalConfigs, response.Luchador.Configs)
	assert.Greater(t, changed, 0)

	log.WithFields(log.Fields{
		"changed": changed,
	}).Info("comparing response.Configs with original.Configs")

	// check if configs are updated in the subsequent GET of luchador
	afterUpdateLuchador := GetLuchador(t, session)
	assert.Equal(t, len(randomConfigs), len(afterUpdateLuchador.Configs))
	AssertConfigMatch(t, randomConfigs, afterUpdateLuchador.Configs)
	changed = CountChangesConfigMatch(t, afterUpdateLuchador.Configs, originalConfigs)
	assert.Greater(t, changed, 0)

	// check if after update the correct event is published
	channel := fmt.Sprintf("luchador.%v.update", luchador.ID)

	log.WithFields(log.Fields{
		"expected":         channel,
		"publishedChannel": mockPublisher.LastChannel,
	}).Info("publish event")
	assert.True(t, mockPublisher.LastChannel == channel)

}
