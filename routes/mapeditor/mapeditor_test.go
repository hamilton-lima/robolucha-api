package mapeditor_test

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

func Setup(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	os.Setenv("GIN_MODE", "release")

	os.Remove(test.DB_NAME)
	ds = datasource.NewDataSource(datasource.BuildSQLLiteConfig(test.DB_NAME))

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher
}

func TestEmptyResult(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	// creates a system game definition
	gd := model.BuildDefaultGameDefinition()
	gd.Name = "TestPlayRequestHandler"
	ds.CreateGameDefinition(&gd)

	// no gamedefinition should return
	gameDefinitions := ds.FindGameDefinitionByOwner(1)
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
	gameDefinitions := ds.FindGameDefinitionByOwner(2)
	result := *gameDefinitions
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Name, "MY-GAME")
}
