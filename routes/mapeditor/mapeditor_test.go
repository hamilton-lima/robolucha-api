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

	err := handler.Update(1, &other)
	assert.True(t, err != nil)
}

func TestUpdateNewOne(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "Me AGAIN"

	err := handler.Update(1, &gd)
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

	err := handler.Update(1, &gd)
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
	err = handler.Update(1, &result[0])
	assert.True(t, err == nil)

	gameDefinitions = handler.Find(1)
	result2 := *gameDefinitions
	assert.Equal(t, len(result2), 1)
	assert.Equal(t, result2[0].ID, result[0].ID)
	assert.Equal(t, "SOME OTHER(2)", result2[0].Name)
}
