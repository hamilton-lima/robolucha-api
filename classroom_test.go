package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/test"
)

func SetupClassroom(t *testing.T) {
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
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher

	router = createRouter(test.API_KEY, "true", SessionAllwaysValid)
}

func TestAddClassroom(t *testing.T) {
	SetupClassroom(t)
	defer dataSource.db.Close()
	classroom := model.Classroom{Name: "testClassroom"}

	plan, _ := json.Marshal(classroom)
	body := string(plan)

	log.WithFields(log.Fields{
		"classroom": classroom.Name,
	}).Debug("classroom before save")

	w := test.PerformRequestNoAuth(router, "POST", "/private/classroom", body)
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Classroom
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response": response,
		"body":     string(w.Body.Bytes()),
	}).Debug("after create")

	assert.Equal(t, response.Name, "testClassroom")
	assert.True(t, len(response.AccessCode) > 0)
	assert.True(t, len(response.Students) == 0)
	assert.True(t, response.OwnerID == 1)

}

func AddTestClassroom(t *testing.T, name string) {
	classroom := model.Classroom{Name: name}
	plan, _ := json.Marshal(classroom)
	body := string(plan)

	log.WithFields(log.Fields{
		"classroom": classroom.Name,
	}).Debug("classroom before save")

	w := test.PerformRequestNoAuth(router, "POST", "/private/classroom", body)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetClassroom(t *testing.T) {
	SetupClassroom(t)
	defer dataSource.db.Close()

	AddTestClassroom(t, "A")
	AddTestClassroom(t, "B")
	AddTestClassroom(t, "C")

	w := test.PerformRequestNoAuth(router, "GET", "/private/classroom", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.Classroom
	json.Unmarshal(w.Body.Bytes(), &response)

	log.WithFields(log.Fields{
		"response": response,
	}).Debug("after create")

	assert.Equal(t, response[0].Name, "A")
	assert.Equal(t, response[1].Name, "B")
	assert.Equal(t, response[2].Name, "C")

	for _, classroom := range response {
		assert.True(t, len(classroom.AccessCode) > 0)
		assert.True(t, len(classroom.Students) == 0)
		assert.True(t, classroom.OwnerID == 1)
	}

}

func TestJoinClassroom(t *testing.T) {
	SetupClassroom(t)
	defer dataSource.db.Close()

	AddTestClassroom(t, "A")
	AddTestClassroom(t, "B")
	AddTestClassroom(t, "C")

	w := test.PerformRequestNoAuth(router, "GET", "/private/classroom", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.Classroom
	json.Unmarshal(w.Body.Bytes(), &response)

	url := fmt.Sprintf("/private/join-classroom/%v", response[1].AccessCode)

	w = test.PerformRequestNoAuth(router, "POST", url, "")
	assert.Equal(t, http.StatusOK, w.Code)
	var joinedClassroom model.Classroom
	json.Unmarshal(w.Body.Bytes(), &joinedClassroom)

	assert.Equal(t, joinedClassroom.Name, "B")

	w = test.PerformRequestNoAuth(router, "GET", "/private/classroom", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var afterJoined []model.Classroom
	json.Unmarshal(w.Body.Bytes(), &afterJoined)

	assert.True(t, len(afterJoined[0].Students) == 0)
	assert.True(t, len(afterJoined[1].Students) == 1)
	assert.True(t, len(afterJoined[2].Students) == 0)
	assert.True(t, afterJoined[1].Students[0].UserID == 1)

}
