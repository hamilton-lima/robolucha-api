package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/bxcodec/faker"
	// "github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/test"
	"gotest.tools/assert"
)

const TEST_USERNAME = "foo"

func SetupMain(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	os.Setenv("GIN_MODE", "release")
}

func TestCreateMatch(t *testing.T) {
	SetupMain(t)
	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	gd := BuildDefaultGameDefinition()
	gd.Name = "FOOBAR"
	dataSource.createGameDefinition(&gd)

	url := fmt.Sprintf("/internal/start-match/%v", gd.Name)

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestCreateTutorialMatch(t *testing.T) {
	SetupMain(t)
	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher

	gd := BuildDefaultGameDefinition()
	gd.Name = "FOOBAR"
	dataSource.createGameDefinition(&gd)

	url := fmt.Sprintf("/private/start-tutorial-match/%v", gd.Name)
	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)

	// force the luchador creation for the current user
	luchadorResponse := test.PerformRequestNoAuth(router, "GET", "/private/luchador", "")
	assert.Equal(t, http.StatusOK, luchadorResponse.Code)
	var luchador GameComponent
	json.Unmarshal(luchadorResponse.Body.Bytes(), &luchador)

	w := test.PerformRequest(router, "POST", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var apiResult *JoinMatch
	json.Unmarshal(w.Body.Bytes(), &apiResult)

	var publisherResult *JoinMatch
	json.Unmarshal([]byte(mockPublisher.LastMessage), &publisherResult)

	match := (*dataSource.findActiveMatches())[0]

	assert.Equal(t, match.ID, publisherResult.MatchID)
	assert.Equal(t, luchador.ID, publisherResult.LuchadorID)
	assert.Equal(t, match.ID, apiResult.MatchID)
	assert.Equal(t, luchador.ID, apiResult.LuchadorID)

	// add participant
	matchParticipant := MatchParticipant{LuchadorID: luchador.ID, MatchID: match.ID}
	body, _ := json.Marshal(matchParticipant)

	w = test.PerformRequest(router, "POST", "/internal/match-participant", string(body), test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	// no messsage should be sent to the publisher if the match exists
	mockPublisher.LastMessage = "EMPTY"

	//call again and expect the same matchID
	w = test.PerformRequest(router, "POST", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var secondApiResult *JoinMatch
	json.Unmarshal(w.Body.Bytes(), &secondApiResult)

	assert.Equal(t, match.ID, secondApiResult.MatchID)
	assert.Equal(t, luchador.ID, secondApiResult.LuchadorID)

	// even if the match exists should send the message to start
	// the runner should only start ONCE
	var publisherResult2 *JoinMatch
	json.Unmarshal([]byte(mockPublisher.LastMessage), &publisherResult2)
	assert.Equal(t, match.ID, publisherResult2.MatchID)
	assert.Equal(t, luchador.ID, publisherResult2.LuchadorID)

	// end match and call again expecting to have a new match
	match.TimeEnd = time.Now()
	body, _ = json.Marshal(match)
	w = test.PerformRequest(router, "PUT", "/internal/end-match", string(body), test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	w = test.PerformRequest(router, "POST", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var thirdApiResult *JoinMatch
	json.Unmarshal(w.Body.Bytes(), &thirdApiResult)

	var thirdPublisherResult *JoinMatch
	json.Unmarshal([]byte(mockPublisher.LastMessage), &thirdPublisherResult)

	thirdCallMatch := (*dataSource.findActiveMatches())[0]

	log.WithFields(log.Fields{
		"thirdApiResult":       thirdApiResult,
		"thirdPublisherResult": thirdPublisherResult,
		"thirdCallMatch":       thirdCallMatch,
		"step":                 "Third call",
	}).Debug("TestCreateTutorialMatch")

	assert.Equal(t, thirdCallMatch.ID, thirdApiResult.MatchID)
	assert.Equal(t, luchador.ID, thirdApiResult.LuchadorID)

	assert.Equal(t, thirdCallMatch.ID, thirdPublisherResult.MatchID)
	assert.Equal(t, luchador.ID, thirdPublisherResult.LuchadorID)

	// make sure it created a new match
	assert.Assert(t, match.ID != thirdApiResult.MatchID)
	assert.Assert(t, match.ID != thirdPublisherResult.MatchID)

}

func TestUpdateGameDefinition(t *testing.T) {
	SetupMain(t)
	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher

	gd, _, _ := fakeGameDefinition(t, "FOOBAR", "tutorial", 10)
	created := dataSource.createGameDefinition(&gd)

	queryResult := dataSource.findGameDefinitionByName(gd.Name)
	assert.Equal(t, created.ID, queryResult.ID)

	log.WithFields(log.Fields{
		"original":      gd.Name,
		"created.Name":  created.Name,
		"created.ID":    created.ID,
		"query by name": queryResult,
	}).Debug("TestUpdateGameDefinition")

	ID := created.ID
	gd.ID = 0

	gd.MinParticipants = 1
	gd.ArenaHeight = 42

	body, _ := json.Marshal(gd)
	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)

	w := test.PerformRequest(router, "PUT", "/internal/game-definition", string(body), test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)
	var updated GameDefinition
	json.Unmarshal(w.Body.Bytes(), &updated)

	log.WithFields(log.Fields{
		"original": gd,
		"created":  created,
		"updated":  updated,
	}).Debug("TestUpdateGameDefinition")

	assert.Equal(t, uint(1), updated.MinParticipants)
	assert.Equal(t, uint(42), updated.ArenaHeight)
	assert.Equal(t, ID, updated.ID)

	// count elements
	assert.Assert(t, len(updated.Codes) == 2)
	assert.Assert(t, len(updated.LuchadorSuggestedCodes) == 2)

	assert.Assert(t, len(updated.GameComponents) == 2)
	assert.Assert(t, len(updated.GameComponents[0].Codes) == 2)
	assert.Assert(t, len(updated.GameComponents[0].Configs) > 0)

	assert.Assert(t, len(updated.SceneComponents) == 2)
	assert.Assert(t, len(updated.SceneComponents[0].Codes) == 2)
}

func TestCreateGameComponent(t *testing.T) {
	SetupMain(t)
	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	plan, _ := ioutil.ReadFile("test-data/create-gamecomponent1.json")
	body := string(plan)
	log.WithFields(log.Fields{
		"body": body,
	}).Debug("After Create Game component")

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/game-component", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var luchador GameComponent
	json.Unmarshal(w.Body.Bytes(), &luchador)
	assert.Assert(t, luchador.ID > 0)
	log.WithFields(log.Fields{
		"luchador.ID": luchador.ID,
	}).Debug("First call to create game component")

	// retry to check for duplicate
	w = test.PerformRequest(router, "POST", "/internal/game-component", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var luchador2 GameComponent
	json.Unmarshal(w.Body.Bytes(), &luchador2)
	assert.Assert(t, luchador.ID == luchador2.ID)
	log.WithFields(log.Fields{
		"luchador.ID": luchador2.ID,
	}).Debug("Second call to create game component")

	luchadorFromDB := dataSource.findLuchadorByID(luchador2.ID)
	log.WithFields(log.Fields{
		"luchador.configs": luchadorFromDB.Configs,
	}).Debug("configs from luchador")

	// all the Mask config items should be present
	for _, color := range maskColors {
		found := false
		for _, config := range luchadorFromDB.Configs {
			if config.Key == color {
				found = true
				break
			}
		}
		assert.Assert(t, found)
		log.WithFields(log.Fields{
			"color": color,
		}).Debug("Color found in luchador config")
	}

	for shape, _ := range maskShapes {
		found := false
		for _, config := range luchadorFromDB.Configs {
			if config.Key == shape {
				found = true
				break
			}
		}
		assert.Assert(t, found)
		log.WithFields(log.Fields{
			"shape": shape,
		}).Debug("Shape found in luchador config")
	}

	getConfigs(t, router, luchador2.ID)
	configsFromDB := getConfigs(t, router, luchador2.ID)
	AssertConfigMatch(t, luchadorFromDB.Configs, configsFromDB)
}

func TestAddScores(t *testing.T) {
	SetupMain(t)
	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	luchador1 := dataSource.createLuchador(&GameComponent{Name: "foo"})
	luchador2 := dataSource.createLuchador(&GameComponent{Name: "bar"})
	luchador3 := dataSource.createLuchador(&GameComponent{Name: "dee"})

	gd := BuildDefaultGameDefinition()
	dataSource.createGameDefinition(&gd)

	match := dataSource.createMatch(gd.ID)
	dataSource.addMatchParticipant(&MatchParticipant{LuchadorID: luchador1.ID, MatchID: match.ID})
	dataSource.addMatchParticipant(&MatchParticipant{LuchadorID: luchador2.ID, MatchID: match.ID})
	dataSource.addMatchParticipant(&MatchParticipant{LuchadorID: luchador3.ID, MatchID: match.ID})

	matchID := fmt.Sprintf("%v", match.ID)
	luchador1ID := fmt.Sprintf("%v", luchador1.ID)
	luchador2ID := fmt.Sprintf("%v", luchador2.ID)
	luchador3ID := fmt.Sprintf("%v", luchador3.ID)

	plan, _ := ioutil.ReadFile("test-data/add-match-scores.json")
	body := string(plan)

	body = strings.Replace(body, "{{.matchID}}", matchID, -1)
	body = strings.Replace(body, "{{.luchadorID1}}", luchador1ID, -1)
	body = strings.Replace(body, "{{.luchadorID2}}", luchador2ID, -1)
	body = strings.Replace(body, "{{.luchadorID3}}", luchador3ID, -1)

	log.WithFields(log.Fields{
		"body": body,
	}).Debug("TestAddScores")

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/add-match-scores", body, test.API_KEY)
	resultScores := dataSource.getMatchScoresByMatchID(match.ID)
	assert.Equal(t, 3, len(*resultScores))
	assert.Equal(t, http.StatusOK, w.Code)

	// parse request body in object to validate the result
	var scoreList ScoreList
	json.Unmarshal([]byte(body), &scoreList)

	// check if all data was saved correctly
	for _, scoreFromBody := range scoreList.Scores {
		found := false
		for _, scoreFromDB := range *resultScores {
			if scoreFromBody.LuchadorID == scoreFromDB.LuchadorID &&
				scoreFromBody.Kills == scoreFromDB.Kills &&
				scoreFromBody.Deaths == scoreFromDB.Deaths &&
				scoreFromBody.Score == scoreFromDB.Score {

				found = true
				break
			}
		}
		assert.Assert(t, found)
		log.WithFields(log.Fields{
			"score-from-body": scoreFromBody,
		}).Debug("TestAddScores")
	}

}

func getConfigs(t *testing.T, router *gin.Engine, id uint) []Config {

	w := test.PerformRequestNoAuth(router, "GET", fmt.Sprintf("/private/mask-config/%v", id), "")
	assert.Equal(t, http.StatusOK, w.Code)

	var configs []Config
	json.Unmarshal(w.Body.Bytes(), &configs)

	return configs
}

func SessionAllwaysValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := dataSource.createUser(User{Username: TEST_USERNAME})
		c.Set("user", user)
	}
}

func TestCreateGameDefinition(t *testing.T) {
	SetupMain(t)

	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	resultFake, body, err := fakeGameDefinition(t, faker.Word(), faker.Word(), 0)
	assert.Assert(t, err == nil)

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/game-definition", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	resultGameDefinition := GameDefinition{}
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)
	compareGameDefinition(t, resultFake, resultGameDefinition)
}
func TestGETGameDefinition(t *testing.T) {
	SetupMain(t)

	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	definition1 := createTestGameDefinition(t, faker.Word(), 0)
	definition2 := createTestGameDefinition(t, faker.Word(), 0)

	url := fmt.Sprintf("/internal/game-definition/%v", definition1.Name)

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "GET", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	resultGameDefinition := GameDefinition{}
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)

	compareGameDefinition(t, definition1, resultGameDefinition)
	assert.Assert(t, definition1.ID != definition2.ID)
}

func TestFindTutorialGameDefinition(t *testing.T) {
	SetupMain(t)

	os.Remove(test.DB_NAME)
	dataSource = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
	defer dataSource.db.Close()

	definition1 := createTestGameDefinition(t, GAMEDEFINITION_TYPE_TUTORIAL, 20)
	definition2 := createTestGameDefinition(t, GAMEDEFINITION_TYPE_TUTORIAL, 10)
	definition3 := createTestGameDefinition(t, faker.Word(), 1)

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequestNoAuth(router, "GET", "/private/tutorial", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var resultGameDefinition []GameDefinition
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)

	// fmt.Printf(">>>>%v", string(w.Body.Bytes()))
	assert.Assert(t, definition3.ID > 0)
	assert.Assert(t, len(resultGameDefinition) == 2)
	compareGameDefinition(t, definition2, resultGameDefinition[0])
	compareGameDefinition(t, definition1, resultGameDefinition[1])
}

func createTestGameDefinition(t *testing.T, typeName string, sortOrder uint) GameDefinition {
	_, body, err := fakeGameDefinition(t, faker.Word(), typeName, sortOrder)
	assert.Assert(t, err == nil)

	router := createRouter(test.API_KEY, "true", SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/game-definition", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	resultGameDefinition := GameDefinition{}
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)

	return resultGameDefinition
}

// assert.DeepEqual is checking not exported fields
func compareGameDefinition(t *testing.T, a, b GameDefinition) {
	assert.Assert(t, len(a.Codes) > 0)
	assert.Assert(t, len(a.LuchadorSuggestedCodes) > 0)

	assert.Assert(t, len(a.GameComponents) > 0)
	assert.Assert(t, len(a.GameComponents[0].Codes) > 0)
	assert.Assert(t, len(a.GameComponents[0].Configs) > 0)

	assert.Assert(t, len(a.SceneComponents) > 0)
	assert.Assert(t, len(a.SceneComponents[0].Codes) > 0)

	assert.Assert(t, a.Name == b.Name)

	assert.Equal(t, len(a.Codes), len(b.Codes))
	assert.Equal(t, len(a.GameComponents), len(b.GameComponents))

	assert.Equal(t, len(a.GameComponents[0].Codes), len(b.GameComponents[0].Codes))
	assert.Equal(t, len(a.GameComponents[0].Configs), len(b.GameComponents[0].Configs))
}

func fakeGameDefinition(t *testing.T, name string, typeName string, sortOrder uint) (GameDefinition, string, error) {
	gameDefinition := GameDefinition{}

	err := faker.FakeData(&gameDefinition)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error generation fake game definition")
		return GameDefinition{}, "", err
	}

	gameDefinition.ID = 0
	gameDefinition.Type = typeName
	gameDefinition.SortOrder = sortOrder
	gameDefinition.Name = name

	gameDefinition.GameComponents = make([]GameComponent, 2)
	gameDefinition.SceneComponents = make([]SceneComponent, 2)
	gameDefinition.Codes = make([]Code, 2)
	gameDefinition.LuchadorSuggestedCodes = make([]Code, 2)

	for i, _ := range gameDefinition.GameComponents {
		faker.FakeData(&gameDefinition.GameComponents[i])

		gameDefinition.GameComponents[i].Codes = make([]Code, 2)
		for n, _ := range gameDefinition.GameComponents[i].Codes {
			faker.FakeData(&gameDefinition.GameComponents[i].Codes[n])
			// gameDefinition.GameComponents[i].Codes[n].ID = 0
		}

		gameDefinition.GameComponents[i].Configs = randomConfig()
	}

	for i, _ := range gameDefinition.SceneComponents {
		faker.FakeData(&gameDefinition.SceneComponents[i])

		gameDefinition.SceneComponents[i].Codes = make([]Code, 2)
		for n, _ := range gameDefinition.SceneComponents[i].Codes {
			faker.FakeData(&gameDefinition.SceneComponents[i].Codes[n])
			// gameDefinition.SceneComponents[i].Codes[n].ID = 0
		}
	}

	for i, _ := range gameDefinition.Codes {
		faker.FakeData(&gameDefinition.Codes[i])
		// gameDefinition.Codes[i].ID = 0
	}

	for i, _ := range gameDefinition.LuchadorSuggestedCodes {
		faker.FakeData(&gameDefinition.LuchadorSuggestedCodes[i])
		// gameDefinition.LuchadorSuggestedCodes[i].ID = 0
	}

	foo, _ := json.Marshal(gameDefinition)
	result := string(foo)

	// removes dates from generated records
	json.Unmarshal([]byte(result), &gameDefinition)

	assert.Assert(t, len(gameDefinition.Codes) == 2)
	assert.Assert(t, len(gameDefinition.LuchadorSuggestedCodes) == 2)

	assert.Assert(t, len(gameDefinition.GameComponents) == 2)
	assert.Assert(t, len(gameDefinition.GameComponents[0].Codes) == 2)
	assert.Assert(t, len(gameDefinition.GameComponents[0].Configs) > 0)

	assert.Assert(t, len(gameDefinition.SceneComponents) == 2)
	assert.Assert(t, len(gameDefinition.SceneComponents[0].Codes) == 2)

	return gameDefinition, result, nil
}
