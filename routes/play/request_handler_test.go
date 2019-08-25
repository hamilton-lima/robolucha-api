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
var publisher pubsub.Publisher
var router *gin.Engine

func Setup(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	os.Setenv("GIN_MODE", "release")

	os.Remove(test.DB_NAME)
	ds = datasource.NewDataSource(datasource.BuildSQLLiteConfig(test.DB_NAME))

	publisher = &test.MockPublisher{}
}

func TestPlayRequestHandler(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	gd := model.BuildDefaultGameDefinition()
	gd.Name = "FOOBAR"
	ds.CreateGameDefinition(&gd)

	am1 := model.AvailableMatch{ID: 42, GameDefinitionID: gd.ID}
	am3 := model.AvailableMatch{ID: 3, GameDefinitionID: gd.ID}

	handler := play.Listen(ds, publisher)

	s1 := handler.Send(play.Request{AvailableMatch: &am1, LuchadorID: 432})
	s2 := handler.Send(play.Request{AvailableMatch: &am1, LuchadorID: 450})

	r1 := <-s1
	r2 := <-s2

	assert.Equal(t, uint(42), r1.Match.AvailableMatchID)
	assert.Equal(t, uint(42), r2.Match.AvailableMatchID)

	s3 := handler.Send(play.Request{AvailableMatch: &am3, LuchadorID: 777})
	r3 := <-s3
	assert.Equal(t, uint(3), r3.Match.AvailableMatchID)

}
