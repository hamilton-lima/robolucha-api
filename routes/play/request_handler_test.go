package play_test

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
	"gitlab.com/robolucha/robolucha-api/routes/play"
	"gitlab.com/robolucha/robolucha-api/test"
	"gotest.tools/assert"
	"testing"
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

func TestPlayRequestHandler(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	gd := model.BuildDefaultGameDefinition()
	gd.Name = "FOOBAR"
	ds.CreateGameDefinition(&gd)

	am1 := model.AvailableMatch{ID: 42, GameDefinitionID: gd.ID}
	am3 := model.AvailableMatch{ID: 3, GameDefinitionID: gd.ID}

	handler := play.NewRequestHandler(ds, publisher)

	r1 := handler.Play(&am1, 432)
	r2 := handler.Play(&am1, 450)
	r3 := handler.Play(&am1, 450)

	startMatchMessages := mockPublisher.Messages["start.match"]
	joinMatchMessages := mockPublisher.Messages["join.match"]

	assert.Equal(t, len(startMatchMessages), 1)
	assert.Equal(t, len(joinMatchMessages), 3)

	assert.Equal(t, uint(42), r1.AvailableMatchID)
	assert.Equal(t, uint(42), r2.AvailableMatchID)
	assert.Equal(t, uint(42), r3.AvailableMatchID)

	r4 := handler.Play(&am3, 777)
	assert.Equal(t, uint(3), r4.AvailableMatchID)

}
