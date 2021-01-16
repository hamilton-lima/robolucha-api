package mapeditor

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
	"gitlab.com/robolucha/robolucha-api/test"
)

var ds *datasource.DataSource
var mockPublisher *test.MockPublisher
var publisher pubsub.Publisher
var router *gin.Engine
var handler *RequestHandler

func Setup(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	os.Setenv("GIN_MODE", "release")

	os.Remove(test.DB_NAME)
	ds = datasource.NewDataSource(datasource.BuildSQLLiteConfig(test.DB_NAME))

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher
	handler = NewRequestHandler(ds, publisher)
}

func TestEmptyResult(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "TestPlayRequestHandler"
	ds.CreateGameDefinition(&gd)

	// no gamedefinition should return
	gameDefinitions := handler.Find(1)
	assert.Equal(t, len(*gameDefinitions), 0)
}

func TestAssignOwner(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "MY-GAME"
	gd.OwnerUserID = 2
	ds.CreateGameDefinition(&gd)

	// no gamedefinition should return
	gameDefinitions := handler.Find(2)
	result := *gameDefinitions
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "MY-GAME")
}

func TestAddAlreadyExist(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "Me AGAIN"
	ds.CreateGameDefinition(&gd)

	err := handler.Add(1, &gd)
	assert.True(t, err != nil)
}
func TestAdd(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "Me AGAIN"
	ds.CreateGameDefinition(&gd)

	// should add with no issues
	gd.Name = "SOME OTHER"
	err := handler.Add(1, &gd)
	assert.True(t, err == nil)

	// check if ID different
	gameDefinitions := handler.Find(1)
	result := *gameDefinitions
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "SOME OTHER")
	assert.True(t, result[0].ID != gd.ID)
	assert.True(t, result[0].ID != 0)
}

func TestUpdateAlreadyExist(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "Me AGAIN"
	ds.CreateGameDefinition(&gd)

	other := model.BuildDefaultGameDefinition()
	other.Name = "Me AGAIN"

	err := handler.Update(1, &other, false)
	assert.True(t, err != nil)
}

func TestUpdateNewOne(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "Me AGAIN"

	err := handler.Update(1, &gd, false)
	assert.True(t, err != nil)
}

func TestUpdateNotOwner(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "Me AGAIN"
	gd.OwnerUserID = 14
	ds.CreateGameDefinition(&gd)

	err := handler.Update(1, &gd, false)
	assert.True(t, err != nil)
}

func TestUpdateName(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "Me AGAIN"
	ds.CreateGameDefinition(&gd)

	// should add with no issues
	gd.Name = "SOME OTHER"
	err := handler.Add(1, &gd)
	assert.True(t, err == nil)

	// check if ID different
	gameDefinitions := handler.Find(1)
	result := *gameDefinitions
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "SOME OTHER")
	assert.True(t, result[0].ID != gd.ID)
	assert.True(t, result[0].ID != 0)

	// should Update name with no issues
	result[0].Name = "SOME OTHER(2)"
	err = handler.Update(1, &result[0], false)
	assert.True(t, err == nil)

	gameDefinitions = handler.Find(1)
	result2 := *gameDefinitions
	assert.Equal(t, len(result2), 1)
	assert.Equal(t, result2[0].ID, result[0].ID)
	assert.Equal(t, "SOME OTHER(2)", result2[0].Name)
}

func TestGetDefault(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	gd := handler.GetDefault()

	assert.True(t, gd.MinParticipants > 0)
	assert.True(t, gd.MaxParticipants > 0)
	assert.True(t, gd.ArenaWidth > 0)
	assert.True(t, gd.ArenaHeight > 0)
}

func TestUpdateGameComponentCode(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "AGAIN"

	gd.GameComponents = make([]model.GameComponent, 1)
	gd.SceneComponents = make([]model.SceneComponent, 0)
	gd.Codes = make([]model.Code, 0)
	gd.LuchadorSuggestedCodes = make([]model.Code, 0)

	gd.GameComponents[0].Name = "otto"
	gd.GameComponents[0].Configs = make([]model.Config, 0)
	gd.GameComponents[0].Codes = make([]model.Code, 1)
	gd.GameComponents[0].Codes[0] = model.Code{Event: "onStart", Script: "turnGun(90)"}

	err := handler.Add(1, &gd)
	assert.True(t, err == nil)

	// check if ID different
	gameDefinitions := handler.Find(1)
	result := *gameDefinitions
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "AGAIN")
	assert.True(t, result[0].ID != gd.ID)
	assert.True(t, result[0].ID != 0)

	id := result[0].ID

	log.WithFields(log.Fields{
		"ID":                                   id,
		"result[0].GameComponents[0].Codes[0]": result[0].GameComponents[0].Codes[0],
	}).Info("BEFORE UPDATE")

	// update the code
	result[0].GameComponents[0].Codes[0] = model.Code{Event: "all", Script: "--updated"}

	log.WithFields(log.Fields{
		"ID":                                   id,
		"result[0].GameComponents[0].Codes[0]": result[0].GameComponents[0].Codes[0],
	}).Info("AFTER UPDATE")

	err = handler.Update(1, &result[0], false)
	assert.True(t, err == nil)

	gameDefinitions2 := handler.Find(1)
	result2 := *gameDefinitions2

	assert.Equal(t, 1, len(result2))
	assert.Equal(t, id, result2[0].ID)

	// both elements from the list will be present
	assert.Equal(t, 2, len(result2[0].GameComponents[0].Codes))
	code := result2[0].GameComponents[0].Codes[1]
	assert.Equal(t, "all", code.Event)
	assert.Equal(t, "--updated", code.Script)
}
