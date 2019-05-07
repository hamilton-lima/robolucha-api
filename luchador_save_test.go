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

	getLuchador := test.PerformRequest(router, "GET", "/private/luchador", "", session)
	assert.Equal(t, http.StatusOK, getLuchador.Code)
	var luchador Luchador
	json.Unmarshal(getLuchador.Body.Bytes(), &luchador)

	return &luchador
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

	// then try a too large name
	luchador.Name = "123456789012345678901234567890aaaaaa"

	plan2, _ := json.Marshal(luchador)
	body2 := string(plan2)

	log.WithFields(log.Fields{
		"luchador": luchador.Name,
	}).Info("luchador before large update")

	w := test.PerformRequest(router, "PUT", "/private/luchador", body2, session)
	assert.Equal(t, http.StatusOK, w.Code)

	var response UpdateLuchadorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response":        response.Luchador.Name,
		"response.errors": response.Errors,
	}).Info("after luchador update")

	t.Log(response.Errors)
	assert.Greater(t, len(response.Errors), 0)
}

func TestLuchadorUpdateName(t *testing.T) {
	luchador := Setup(t)
	defer dataSource.db.Close()

	// first try to change to a valid name
	luchador.Name = "lucharito"
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
		"response":          response.Luchador.Name,
		"response.luchador": response.Luchador,
	}).Info("after luchador update")

	assert.Equal(t, "lucharito", response.Luchador.Name)
	assert.Equal(t, 0, len(response.Errors))

	channel := fmt.Sprintf("luchador.%v.update", luchador.ID)

	log.WithFields(log.Fields{
		"expected":         channel,
		"publishedChannel": mockPublisher.LastChannel,
	}).Info("publish event")
	assert.True(t, mockPublisher.LastChannel == channel)

}

func getConfig(configs []Config, key string) Config {
	for _, config := range configs {
		if config.Key == key {
			return config
		}
	}
	return Config{}
}
