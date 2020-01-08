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

	"github.com/bxcodec/faker/v3"
	// "github.com/stretchr/testify/assert"
	"gitlab.com/robolucha/robolucha-api/auth"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/test"
	"gotest.tools/assert"
)

func SetupMain(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	os.Setenv("GIN_MODE", "release")

	os.Remove(test.DB_NAME)
	ds = datasource.NewDataSource(datasource.BuildSQLLiteConfig(test.DB_NAME))
}

// func TestCreateMatch(t *testing.T) {
// 	SetupMain(t)
// 	defer ds.DB.Close()

// 	gd := BuildDefaultGameDefinition()
// 	gd.Name = "FOOBAR"
// 	ds.CreateGameDefinition(&gd)

// 	url := fmt.Sprintf("/internal/start-match/%v", gd.Name)

// 	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid)
// 	w := test.PerformRequest(router, "POST", url, "", test.API_KEY)
// 	assert.Equal(t, http.StatusOK, w.Code)

// }

// func TestCreateTutorialMatch(t *testing.T) {
// 	SetupMain(t)
// 	defer ds.DB.Close()

// 	mockPublisher = &test.MockPublisher{}
// 	publisher = mockPublisher

// 	gd := BuildDefaultGameDefinition()
// 	gd.Name = "FOOBAR"
// 	ds.CreateGameDefinition(&gd)

// 	url := fmt.Sprintf("/private/start-tutorial-match/%v", gd.Name)
// 	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid)

// 	// force the luchador creation for the current user
// 	luchadorResponse := test.PerformRequestNoAuth(router, "GET", "/private/luchador", "")
// 	assert.Equal(t, http.StatusOK, luchadorResponse.Code)
// 	var luchador model.GameComponent
// 	json.Unmarshal(luchadorResponse.Body.Bytes(), &luchador)

// 	w := test.PerformRequest(router, "POST", url, "", test.API_KEY)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var apiResult *model.JoinMatch
// 	json.Unmarshal(w.Body.Bytes(), &apiResult)

// 	var publisherResult *model.JoinMatch
// 	json.Unmarshal([]byte(mockPublisher.LastMessage), &publisherResult)

// 	match := (*ds.FindActiveMatches())[0]

// 	assert.Equal(t, match.ID, publisherResult.MatchID)
// 	assert.Equal(t, luchador.ID, publisherResult.LuchadorID)
// 	assert.Equal(t, match.ID, apiResult.MatchID)
// 	assert.Equal(t, luchador.ID, apiResult.LuchadorID)

// 	// add participant
// 	matchParticipant := model.MatchParticipant{LuchadorID: luchador.ID, MatchID: match.ID}
// 	body, _ := json.Marshal(matchParticipant)

// 	w = test.PerformRequest(router, "POST", "/internal/match-participant", string(body), test.API_KEY)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// no messsage should be sent to the publisher if the match exists
// 	mockPublisher.LastMessage = "EMPTY"

// 	//call again and expect the same matchID
// 	w = test.PerformRequest(router, "POST", url, "", test.API_KEY)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var secondApiResult *model.JoinMatch
// 	json.Unmarshal(w.Body.Bytes(), &secondApiResult)

// 	assert.Equal(t, match.ID, secondApiResult.MatchID)
// 	assert.Equal(t, luchador.ID, secondApiResult.LuchadorID)

// 	// even if the match exists should send the message to start
// 	// the runner should only start ONCE
// 	var publisherResult2 *model.JoinMatch
// 	json.Unmarshal([]byte(mockPublisher.LastMessage), &publisherResult2)
// 	assert.Equal(t, match.ID, publisherResult2.MatchID)
// 	assert.Equal(t, luchador.ID, publisherResult2.LuchadorID)

// 	// end match and call again expecting to have a new match
// 	match.TimeEnd = time.Now()
// 	body, _ = json.Marshal(match)
// 	w = test.PerformRequest(router, "PUT", "/internal/end-match", string(body), test.API_KEY)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	w = test.PerformRequest(router, "POST", url, "", test.API_KEY)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var thirdAPIResult *model.JoinMatch
// 	json.Unmarshal(w.Body.Bytes(), &thirdAPIResult)

// 	var thirdPublisherResult *model.JoinMatch
// 	json.Unmarshal([]byte(mockPublisher.LastMessage), &thirdPublisherResult)

// 	thirdCallMatch := (*ds.FindActiveMatches())[0]

// 	log.WithFields(log.Fields{
// 		"thirdApiResult":       thirdAPIResult,
// 		"thirdPublisherResult": thirdPublisherResult,
// 		"thirdCallMatch":       thirdCallMatch,
// 		"step":                 "Third call",
// 	}).Debug("TestCreateTutorialMatch")

// 	assert.Equal(t, thirdCallMatch.ID, thirdAPIResult.MatchID)
// 	assert.Equal(t, luchador.ID, thirdAPIResult.LuchadorID)

// 	assert.Equal(t, thirdCallMatch.ID, thirdPublisherResult.MatchID)
// 	assert.Equal(t, luchador.ID, thirdPublisherResult.LuchadorID)

// 	// make sure it created a new match
// 	assert.Assert(t, match.ID != thirdAPIResult.MatchID)
// 	assert.Assert(t, match.ID != thirdPublisherResult.MatchID)

// }

func TestUpdateGameDefinition(t *testing.T) {
	SetupMain(t)
	defer ds.DB.Close()

	mockPublisher = &test.MockPublisher{}
	publisher = mockPublisher

	gd, _, _ := fakeGameDefinition(t, "FOOBAR", "tutorial", 10)
	created := ds.CreateGameDefinition(&gd)

	queryResult := ds.FindGameDefinitionByName(gd.Name)
	assert.Equal(t, created.ID, queryResult.ID)

	log.WithFields(log.Fields{
		"original":      gd.Name,
		"created.Name":  created.Name,
		"created.ID":    created.ID,
		"query by name": queryResult,
	}).Debug("TestUpdateGameDefinition")

	ID := created.ID
	gd.ID = 0

	// remove IDS from GameComponents to create with a differnt name
	gd.GameComponents[0].ID = 0
	gd.GameComponents[0].GameDefinitionID = 0
	gd.GameComponents[0].Name = gd.GameComponents[0].Name + "-UPDATED"

	gd.GameComponents[1].ID = 0
	gd.GameComponents[1].GameDefinitionID = 0

	log.WithFields(log.Fields{
		"gd.GameComponents": gd.GameComponents,
	}).Debug("Before update")

	gd.MinParticipants = 1
	gd.ArenaHeight = 42

	body, _ := json.Marshal(gd)
	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)

	w := test.PerformRequest(router, "PUT", "/internal/game-definition", string(body), test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)
	var updated model.GameDefinition
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

	counterUpdatedCounter := 0
	for _, gc := range updated.GameComponents {

		log.WithFields(log.Fields{
			"gc.Name": gc.Name,
		}).Debug("counterUpdatedCounter")

		if strings.HasSuffix(gc.Name, "-UPDATED") {
			counterUpdatedCounter = counterUpdatedCounter + 1
		}
	}
	assert.Assert(t, counterUpdatedCounter == 1)

	url := fmt.Sprintf("/internal/game-definition/%v", updated.Name)
	w = test.PerformRequest(router, "GET", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	resultGameDefinition := model.GameDefinition{}
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)

	assert.Assert(t, len(resultGameDefinition.GameComponents) == 2)

	counterUpdatedCounter2 := 0
	for _, gc := range resultGameDefinition.GameComponents {
		log.WithFields(log.Fields{
			"gc.Name": gc.Name,
		}).Debug("counterUpdatedCounter2")

		if strings.HasSuffix(gc.Name, "-UPDATED") {
			counterUpdatedCounter2 = counterUpdatedCounter2 + 1
		}
	}
	assert.Assert(t, counterUpdatedCounter2 == 1)
}

func TestCreateGameComponent(t *testing.T) {
	SetupMain(t)
	defer ds.DB.Close()

	plan, _ := ioutil.ReadFile("test-data/create-gamecomponent1.json")
	body := string(plan)
	log.WithFields(log.Fields{
		"body": body,
	}).Debug("After Create Game component")

	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/game-component", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var luchador model.GameComponent
	json.Unmarshal(w.Body.Bytes(), &luchador)
	assert.Assert(t, luchador.ID > 0)
	log.WithFields(log.Fields{
		"luchador.ID": luchador.ID,
	}).Debug("First call to create game component")

	// retry to check for duplicate
	w = test.PerformRequest(router, "POST", "/internal/game-component", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var luchador2 model.GameComponent
	json.Unmarshal(w.Body.Bytes(), &luchador2)
	assert.Assert(t, luchador.ID == luchador2.ID)
	log.WithFields(log.Fields{
		"luchador.ID": luchador2.ID,
	}).Debug("Second call to create game component")

	luchadorFromDB := ds.FindLuchadorByID(luchador2.ID)
	log.WithFields(log.Fields{
		"luchador.configs": luchadorFromDB.Configs,
	}).Debug("configs from luchador")

	// all the Mask config items should be present
	for _, color := range model.MaskColors {
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

	for shape, _ := range model.MaskShapes {
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

// func TestAddScores(t *testing.T) {
// 	SetupMain(t)
// 	defer ds.DB.Close()

// 	luchador1 := ds.CreateLuchador(&model.GameComponent{Name: "foo"})
// 	luchador2 := ds.CreateLuchador(&model.GameComponent{Name: "bar"})
// 	luchador3 := ds.CreateLuchador(&model.GameComponent{Name: "dee"})

// 	gd := BuildDefaultGameDefinition()
// 	ds.CreateGameDefinition(&gd)

// 	match := ds.CreateMatch(gd.ID)
// 	ds.AddMatchParticipant(&model.MatchParticipant{LuchadorID: luchador1.ID, MatchID: match.ID})
// 	ds.AddMatchParticipant(&model.MatchParticipant{LuchadorID: luchador2.ID, MatchID: match.ID})
// 	ds.AddMatchParticipant(&model.MatchParticipant{LuchadorID: luchador3.ID, MatchID: match.ID})

// 	matchID := fmt.Sprintf("%v", match.ID)
// 	luchador1ID := fmt.Sprintf("%v", luchador1.ID)
// 	luchador2ID := fmt.Sprintf("%v", luchador2.ID)
// 	luchador3ID := fmt.Sprintf("%v", luchador3.ID)

// 	plan, _ := ioutil.ReadFile("test-data/add-match-scores.json")
// 	body := string(plan)

// 	body = strings.Replace(body, "{{.matchID}}", matchID, -1)
// 	body = strings.Replace(body, "{{.luchadorID1}}", luchador1ID, -1)
// 	body = strings.Replace(body, "{{.luchadorID2}}", luchador2ID, -1)
// 	body = strings.Replace(body, "{{.luchadorID3}}", luchador3ID, -1)

// 	log.WithFields(log.Fields{
// 		"body": body,
// 	}).Debug("TestAddScores")

// 	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid)
// 	w := test.PerformRequest(router, "POST", "/internal/add-match-scores", body, test.API_KEY)
// 	resultScores := ds.GetMatchScoresByMatchID(match.ID)
// 	assert.Equal(t, 3, len(*resultScores))
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// parse request body in object to validate the result
// 	var scoreList model.ScoreList
// 	json.Unmarshal([]byte(body), &scoreList)

// 	// check if all data was saved correctly
// 	for _, scoreFromBody := range scoreList.Scores {
// 		found := false
// 		for _, scoreFromDB := range *resultScores {
// 			if scoreFromBody.LuchadorID == scoreFromDB.LuchadorID &&
// 				scoreFromBody.Kills == scoreFromDB.Kills &&
// 				scoreFromBody.Deaths == scoreFromDB.Deaths &&
// 				scoreFromBody.Score == scoreFromDB.Score {

// 				found = true
// 				break
// 			}
// 		}
// 		assert.Assert(t, found)
// 		log.WithFields(log.Fields{
// 			"score-from-body": scoreFromBody,
// 		}).Debug("TestAddScores")
// 	}

// }

func getConfigs(t *testing.T, router *gin.Engine, id uint) []model.Config {

	w := test.PerformRequestNoAuth(router, "GET", fmt.Sprintf("/private/mask-config/%v", id), "")
	assert.Equal(t, http.StatusOK, w.Code)

	var configs []model.Config
	json.Unmarshal(w.Body.Bytes(), &configs)

	return configs
}

func TestCreateGameDefinition(t *testing.T) {
	SetupMain(t)
	defer ds.DB.Close()

	resultFake, body, err := fakeGameDefinition(t, faker.Word(), faker.Word(), 0)
	assert.Assert(t, err == nil)

	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/game-definition", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	resultGameDefinition := model.GameDefinition{}
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)
	compareGameDefinition(t, resultFake, resultGameDefinition)
}
func TestGETGameDefinition(t *testing.T) {
	SetupMain(t)
	defer ds.DB.Close()

	definition1 := createTestGameDefinition(t, faker.Word(), 0)
	definition2 := createTestGameDefinition(t, faker.Word(), 0)

	url := fmt.Sprintf("/internal/game-definition/%v", definition1.Name)

	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequest(router, "GET", url, "", test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	resultGameDefinition := model.GameDefinition{}
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)

	compareGameDefinition(t, definition1, resultGameDefinition)
	assert.Assert(t, definition1.ID != definition2.ID)
}

func createMatch(gameDefinitionID uint) model.Match {
	match := model.Match{
		TimeStart:        time.Now(),
		GameDefinitionID: gameDefinitionID}

	ds.DB.Create(&match)
	return match
}

func TestFindMultiplayerMatch(t *testing.T) {
	SetupMain(t)
	// log.SetLevel(log.InfoLevel)
	// ds.DB.LogMode(true)
	defer ds.DB.Close()

	definition1 := createTestGameDefinition(t, model.GAMEDEFINITION_TYPE_TUTORIAL, 20)
	definition2 := createTestGameDefinition(t, model.GAMEDEFINITION_TYPE_MULTIPLAYER, 10)
	definition3 := createTestGameDefinition(t, faker.Word(), 1)

	createMatch(definition1.ID)
	match := createMatch(definition2.ID)
	createMatch(definition3.ID)

	definition2.Duration = 1200000
	ds.DB.Save(&definition2)

	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequestNoAuth(router, "GET", "/private/match", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var matches []model.ActiveMatch

	log.WithFields(log.Fields{
		"matches": matches,
	}).Warning("TestFindMultiplayerMatch")

	json.Unmarshal(w.Body.Bytes(), &matches)

	assert.Assert(t, matches[0].MatchID == match.ID)

	gameDefinitions := *ds.FindTutorialGameDefinition()

	// all the tutorial gamedefinitions and the active multiplayer matches
	assert.Assert(t, len(matches) == len(gameDefinitions)+1)
}

func TestFindTutorialGameDefinition(t *testing.T) {
	SetupMain(t)
	defer ds.DB.Close()

	definition1 := createTestGameDefinition(t, model.GAMEDEFINITION_TYPE_TUTORIAL, 20)
	definition2 := createTestGameDefinition(t, model.GAMEDEFINITION_TYPE_TUTORIAL, 10)
	definition3 := createTestGameDefinition(t, faker.Word(), 1)

	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequestNoAuth(router, "GET", "/private/tutorial", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var resultGameDefinition []model.GameDefinition
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)

	// fmt.Printf(">>>>%v", string(w.Body.Bytes()))
	assert.Assert(t, definition3.ID > 0)
	assert.Assert(t, len(resultGameDefinition) == 2)
	compareGameDefinition(t, definition2, resultGameDefinition[0])
	compareGameDefinition(t, definition1, resultGameDefinition[1])
}

func createTestGameDefinition(t *testing.T, typeName string, sortOrder uint) model.GameDefinition {
	_, body, err := fakeGameDefinition(t, faker.Word(), typeName, sortOrder)
	assert.Assert(t, err == nil)

	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/game-definition", body, test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	resultGameDefinition := model.GameDefinition{}
	json.Unmarshal(w.Body.Bytes(), &resultGameDefinition)

	return resultGameDefinition
}

// assert.DeepEqual is checking not exported fields
func compareGameDefinition(t *testing.T, a, b model.GameDefinition) {
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

func fakeGameDefinition(t *testing.T, name string, typeName string, sortOrder uint) (model.GameDefinition, string, error) {
	gameDefinition := model.GameDefinition{}

	err := faker.FakeData(&gameDefinition)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error generation fake game definition")
		return model.GameDefinition{}, "", err
	}

	gameDefinition.ID = 0
	gameDefinition.Type = typeName
	gameDefinition.SortOrder = sortOrder
	gameDefinition.Name = name

	gameDefinition.GameComponents = make([]model.GameComponent, 2)
	gameDefinition.SceneComponents = make([]model.SceneComponent, 2)
	gameDefinition.Codes = make([]model.Code, 2)
	gameDefinition.LuchadorSuggestedCodes = make([]model.Code, 2)

	for i, _ := range gameDefinition.GameComponents {
		faker.FakeData(&gameDefinition.GameComponents[i])

		gameDefinition.GameComponents[i].Codes = make([]model.Code, 2)
		for n, _ := range gameDefinition.GameComponents[i].Codes {
			faker.FakeData(&gameDefinition.GameComponents[i].Codes[n])
			// gameDefinition.GameComponents[i].Codes[n].ID = 0
		}

		gameDefinition.GameComponents[i].Configs = model.RandomConfig()
	}

	for i, _ := range gameDefinition.SceneComponents {
		faker.FakeData(&gameDefinition.SceneComponents[i])

		gameDefinition.SceneComponents[i].Codes = make([]model.Code, 2)
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

func TestPOSTMatchMetric(t *testing.T) {
	SetupMain(t)
	defer ds.DB.Close()

	var metric *model.MatchMetric
	err := faker.FakeData(&metric)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error generation fake match metric")
		panic(1)
	}
	body, _ := json.Marshal(metric)

	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	w := test.PerformRequest(router, "POST", "/internal/match-metric", string(body), test.API_KEY)
	assert.Equal(t, http.StatusOK, w.Code)

	var afterSave model.MatchMetric
	ds.DB.First(&afterSave)

	assert.Assert(t, &afterSave != nil)

	assert.Equal(t, metric.MatchID, afterSave.MatchID)
	assert.Equal(t, metric.FPS, afterSave.FPS)
	assert.Equal(t, metric.Players, afterSave.Players)
	assert.Equal(t, metric.GameDefinitionID, afterSave.GameDefinitionID)
}

func TestGetPublicAvailableMatch(t *testing.T) {
	SetupMain(t)
	defer ds.DB.Close()

	// create classroom
	classroom := model.Classroom{
		Name: faker.Word(),
	}
	ds.DB.Create(&classroom)

	// create activeMatch related to the classroom
	availableMatch := model.AvailableMatch{
		Name:        faker.Word(),
		ClassroomID: classroom.ID,
	}
	ds.DB.Create(&availableMatch)

	// create activeMatch without classroom
	publicAvailableMatch := model.AvailableMatch{
		Name: faker.Word(),
	}
	ds.DB.Create(&publicAvailableMatch)

	router := createRouter(test.API_KEY, "true", auth.SessionAllwaysValid, auth.SessionAllwaysValid)

	// search public available matches
	w := test.PerformRequestNoAuth(router, "GET", "/private/available-match-public", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var result []model.AvailableMatch
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.Assert(t, len(result) == 1)
	assert.Assert(t, result[0].ID == publicAvailableMatch.ID)
	assert.Assert(t, result[0].Name == publicAvailableMatch.Name)

	// search  available matches by classroom
	url := fmt.Sprintf("/private/available-match-classroom/%v", classroom.ID)
	w = test.PerformRequestNoAuth(router, "GET", url, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var result2 []model.AvailableMatch
	json.Unmarshal(w.Body.Bytes(), &result2)

	assert.Assert(t, len(result2) == 1)
	assert.Assert(t, result2[0].ID == availableMatch.ID)
	assert.Assert(t, result2[0].Name == availableMatch.Name)

}

func TestAddCodeVersion(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	code := model.Code{Event: "onStart", Script: "turnGun(90)"}

	ds.DB.Create(&code)
	assert.Equal(t, code.Version, uint(1))

	// search version
	var version model.CodeHistory
	ds.DB.Where(&model.CodeHistory{CodeID: code.ID, Version: code.Version}).First(&version)
	assert.Equal(t, version.Script, "turnGun(90)")

	code.Script = "turnGun(160)"
	ds.DB.Save(&code)
	assert.Equal(t, code.Version, uint(2))

	// search second version
	var second model.CodeHistory
	ds.DB.Where(&model.CodeHistory{CodeID: code.ID, Version: code.Version}).First(&second)
	assert.Equal(t, second.Script, "turnGun(160)")
}

func TestUpdateLuchadorCode(t *testing.T) {
	Setup(t)
	defer ds.DB.Close()

	luchador := model.GameComponent{}
	err := faker.FakeData(&luchador)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error generating mock luchador")
		panic(1)
	}

	luchador.Codes = make([]model.Code, 0)
	ds.UpdateLuchador(&luchador)
	result := ds.FindLuchadorByID(luchador.ID)
	assert.Equal(t, len(result.Codes), 0)

	luchador.Codes = make([]model.Code, 2)
	luchador.Codes[0] = model.Code{Event: "onStart", Script: "turnGun(90)"}
	luchador.Codes[1] = model.Code{Event: "onRepeat", Script: "move(10)"}

	ds.UpdateLuchador(&luchador)
	assert.Equal(t, len(luchador.Codes), 2)
	assert.Equal(t, luchador.Codes[0].Script, "turnGun(90)")

	result = ds.FindLuchadorByID(luchador.ID)
	assert.Equal(t, len(result.Codes), 2)
	assert.Equal(t, result.Codes[0].Script, "turnGun(90)")
	assert.Equal(t, result.Codes[1].Script, "move(10)")

	log.WithFields(log.Fields{
		"luchador.Codes": luchador.Codes,
	}).Error("update (1)")

	// force new objects to create new codes instead of updating
	luchador.Codes[0] = model.Code{Event: "onStart", Script: ""}
	luchador.Codes[1] = model.Code{Event: "onRepeat", Script: "move(99)"}

	ds.UpdateLuchador(&luchador)

	result = ds.FindLuchadorByID(luchador.ID)
	assert.Equal(t, len(result.Codes), 2)
	assert.Equal(t, result.Codes[0].Script, "")
	assert.Equal(t, result.Codes[1].Script, "move(99)")

	log.WithFields(log.Fields{
		"result.Codes": result.Codes,
	}).Error("update (2)")

}
