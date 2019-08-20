package play_test

import (
	"os"

	log "github.com/sirupsen/logrus"

	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
	"gitlab.com/robolucha/robolucha-api/routes/play"
	"gitlab.com/robolucha/robolucha-api/test"
	"gotest.tools/assert"
	"testing"
)

var mockPublisher *test.MockPublisher
var ds *datasource.DataSource
var publisher pubsub.Publisher

func Setup(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	os.Setenv("GIN_MODE", "release")

	os.Remove(test.DB_NAME)
	ds = datasource.NewDataSource(datasource.BuildSQLLiteConfig(test.DB_NAME))

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher
}

func TestPlayRequestHandler(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	handler := play.Listen(ds, publisher)

	s1 := handler.Send(play.Request{AvailableMatch: &model.AvailableMatch{ID: 42}})
	s2 := handler.Send(play.Request{AvailableMatch: &model.AvailableMatch{ID: 43}})

	r1 := <-s1
	r2 := <-s2

	assert.Equal(t, uint(42), r1.Match.ID)
	assert.Equal(t, uint(43), r2.Match.ID)

	s3 := handler.Send(play.Request{AvailableMatch: &model.AvailableMatch{ID: 3}})
	r3 := <-s3
	assert.Equal(t, uint(3), r3.Match.ID)

}
