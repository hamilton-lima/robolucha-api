package play_test

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
	"gitlab.com/robolucha/robolucha-api/routes/play"
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

func TestPlayRequestHandler(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	gd := model.BuildDefaultGameDefinition()
	gd.Name = "TestPlayRequestHandler"
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

func createLuchador(id uint) *model.GameComponent {
	luchador := &model.GameComponent{
		UserID: id,
		Name:   fmt.Sprintf("Luchador%d", id),
	}

	return ds.CreateLuchador(luchador)
}

func TestLeaveTutorial(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()
	handler := play.NewRequestHandler(ds, publisher)

	luchador := createLuchador(1)

	gd := model.BuildDefaultGameDefinition()
	gd.Name = "TestLeaveTutorial"
	gd.Type = model.GAMEDEFINITION_TYPE_TUTORIAL
	gd.Duration = 0
	gdCreated := ds.CreateGameDefinition(&gd)

	am1 := model.AvailableMatch{ID: 42, GameDefinitionID: gdCreated.ID}
	am2 := model.AvailableMatch{ID: 3, GameDefinitionID: gdCreated.ID}

	match := handler.Play(&am1, luchador.ID)
	handler.Play(&am2, 450)
	handler.Play(&am1, 450)

	// simulate runner adding the match participant
	// From Runner: MatchRunnerAPI.getInstance().addMatchParticipant
	ds.AddMatchParticipant(&model.MatchParticipant{
		LuchadorID: luchador.ID,
		MatchID:    match.ID,
	})

	handler.LeaveTutorialMatches(luchador)
	endMatchMessages := mockPublisher.Messages["end.match"]

	log.WithFields(log.Fields{
		"messages": mockPublisher.Messages,
	}).Info("TestLeaveTutorial")

	assert.Equal(t, len(endMatchMessages), 1)
}

func TestUserHasLevelToPlay(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	handler := play.NewRequestHandler(ds, publisher)

	levelZero := model.UserLevel{Level: 0}
	levelTen := model.UserLevel{Level: 10}
	levelTwelve := model.UserLevel{Level: 12}
	levelFourTeen := model.UserLevel{Level: 14}

	defZeroZero := model.GameDefinition{MinLevel: 0, MaxLevel: 0}
	defTenZero := model.GameDefinition{MinLevel: 10, MaxLevel: 0}
	defTenThirteen := model.GameDefinition{MinLevel: 10, MaxLevel: 13}

	assert.True(t, handler.UserHasLevelToPlay(&levelZero, &defZeroZero))
	assert.False(t, handler.UserHasLevelToPlay(&levelZero, &defTenZero))
	assert.False(t, handler.UserHasLevelToPlay(&levelZero, &defTenThirteen))

	assert.True(t, handler.UserHasLevelToPlay(&levelTen, &defZeroZero))
	assert.True(t, handler.UserHasLevelToPlay(&levelTen, &defTenZero))
	assert.True(t, handler.UserHasLevelToPlay(&levelTen, &defTenThirteen))

	assert.True(t, handler.UserHasLevelToPlay(&levelTwelve, &defZeroZero))
	assert.True(t, handler.UserHasLevelToPlay(&levelTwelve, &defTenZero))
	assert.True(t, handler.UserHasLevelToPlay(&levelTwelve, &defTenThirteen))

	assert.True(t, handler.UserHasLevelToPlay(&levelFourTeen, &defZeroZero))
	assert.True(t, handler.UserHasLevelToPlay(&levelFourTeen, &defTenZero))
	assert.False(t, handler.UserHasLevelToPlay(&levelFourTeen, &defTenThirteen))

}
