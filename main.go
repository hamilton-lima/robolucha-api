// @title Robolucha API
// @version 1.0
// @description Robolucha API
// @host http://local.robolucha.com:5000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	useragent "github.com/mileusna/useragent"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gitlab.com/robolucha/robolucha-api/auth"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/events"
	"gitlab.com/robolucha/robolucha-api/httphelper"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
	"gitlab.com/robolucha/robolucha-api/routes"
	"gitlab.com/robolucha/robolucha-api/routes/learning"
	"gitlab.com/robolucha/robolucha-api/routes/mapeditor"
	"gitlab.com/robolucha/robolucha-api/routes/play"
	"gitlab.com/robolucha/robolucha-api/setup"
	"gitlab.com/robolucha/robolucha-api/utility"

	_ "gitlab.com/robolucha/robolucha-api/docs"
)

var ds *datasource.DataSource
var eventsDS *events.DataSource

var publisher pubsub.Publisher

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	rand.Seed(time.Now().UTC().UnixNano())

	log.Info("Robolucha API, start.")

	ds = datasource.NewDataSource(datasource.BuildMysqlConfig())
	defer ds.DB.Close()

	eventsDS = events.NewDataSource(events.BuildMysqlConfig())
	defer eventsDS.DB.Close()

	publisher = &pubsub.RedisPublisher{}
	go ds.KeepAlive()
	go eventsDS.KeepAlive()

	if len(os.Args) < 2 {
		log.Error("Wrong number of parameters, usage: api <metadata folder>")
		os.Exit(2)
	}

	metadataFolder := os.Args[1]
	setup.LoadMetadataFromFolder(metadataFolder, ds)
	setup.CreateAvailableMatches(ds)

	port := os.Getenv("API_PORT")
	if len(port) == 0 {
		port = "5000"
	}

	log.WithFields(log.Fields{
		"port": port,
	}).Debug("Port configuration")

	internalAPIKey := os.Getenv("INTERNAL_API_KEY")
	logRequestBody := os.Getenv("GIM_LOG_REQUEST_BODY")
	disableAuth := os.Getenv("DISABLE_AUTH")

	var router *gin.Engine

	if disableAuth == "true" {
		router = createRouter(internalAPIKey, logRequestBody, auth.SessionAllwaysValid, auth.SessionAllwaysValid)
	} else {
		router = createRouter(internalAPIKey, logRequestBody, auth.SessionIsValid, auth.SessionIsValidAndDashBoardUser)
	}

	router.Run(":" + port)

	log.WithFields(log.Fields{
		"port": port,
	}).Debug("Server is ready")
}

func createRouter(internalAPIKey string, logRequestBody string,
	privateFactory auth.SessionValidatorFactory,
	dashboardFactory auth.SessionValidatorFactory) *gin.Engine {

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))
	if logRequestBody == "true" {
		router.Use(RequestLogger())
	}

	publicAPI := router.Group("/public")
	{
		publicAPI.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	internalAPI := router.Group("/internal")
	internalAPI.Use(auth.KeyIsValid(internalAPIKey))
	{
		internalAPI.GET("/game-definition/:name", getGameDefinitionByName)
		internalAPI.GET("/game-definition-id/:id", getGameDefinitionByIDInternal)
		internalAPI.POST("/game-definition", createGameDefinition)
		internalAPI.PUT("/game-definition", updateGameDefinition)
		internalAPI.POST("/game-component", createGameComponent)
		internalAPI.POST("/luchador", getLuchadorByIDAndGamedefinitionID)
		internalAPI.POST("/match-participant", addMatchPartipant)
		internalAPI.PUT("/end-match", endMatch)
		internalAPI.PUT("/run-match", runMatch)
		internalAPI.GET("/ready", getReady)
		internalAPI.POST("/add-match-scores", addMatchScores)
		internalAPI.GET("/match-single", getMatchInternal)
		internalAPI.POST("/match-metric", addMatchMetric)
	}

	privateAPI := router.Group("/private")
	privateAPI.Use(privateFactory(ds))
	{
		privateAPI.GET("/tutorial", getTutorialGameDefinition)
		privateAPI.GET("/get-user", getUser)
		privateAPI.GET("/luchador", getLuchador)
		privateAPI.PUT("/luchador", updateLuchador)
		privateAPI.GET("/mask-config/:id", getMaskConfig)
		privateAPI.GET("/mask-random", getRandomMaskConfig)
		privateAPI.PUT("/user/setting", updateUserSetting)
		privateAPI.GET("/user/setting", findUserSetting)
		privateAPI.GET("/match", getActiveMatches)
		privateAPI.GET("/match-multiplayer", getActiveMultiplayerMatches)

		privateAPI.GET("/match-single", getMatch)
		privateAPI.GET("/match-score", getMatchScore)
		privateAPI.GET("/match-config", getLuchadorConfigsForCurrentMatch)
		privateAPI.POST("/join-match", joinMatch)
		privateAPI.GET("/game-definition-id/:id", getGameDefinitionByID)
		privateAPI.GET("/game-definition-all", getGameDefinition)
		privateAPI.GET("/classroom", getClassroom)
		privateAPI.POST("/classroom", addClassroom)
		privateAPI.POST("/join-classroom/:accessCode", joinClassroom)
		privateAPI.GET("/available-match-public", getPublicAvailableMatch)
		privateAPI.GET("/available-match-classroom/:id", getClassroomAvailableMatch)
		privateAPI.POST("/page-events", addEvents)
		privateAPI.GET("/level-group", getLevelGroup)

	}

	dashboardAPI := router.Group("/dashboard")
	dashboardAPI.Use(dashboardFactory(ds))
	{
		dashboardAPI.GET("/get-user", getUserDashboard)
		dashboardAPI.GET("/classroom", getClassroom)
		dashboardAPI.POST("/classroom", addClassroom)
		dashboardAPI.GET("/classroom/students/:id", getClassroomStudents)
	}

	learningRouter := learning.Init(ds, publisher)
	routes.Use(dashboardAPI, learningRouter)

	playRouter := play.Init(ds, publisher)
	routes.Use(privateAPI, playRouter)

	mapeditorRouter := mapeditor.Init(ds, publisher)
	routes.Use(privateAPI, mapeditorRouter)

	return router
}

// findUserSetting godoc
// @Summary find current user userSetting
// @Accept  json
// @Produce  json
// @Success 200 {object} model.UserSetting
// @Security ApiKeyAuth
// @Router /private/user/setting [get]
func findUserSetting(c *gin.Context) {

	log.Info("Finding userSetting")
	user := httphelper.UserFromContext(c)

	userSetting := ds.FindUserSettingByUser(user)

	log.WithFields(log.Fields{
		"userSetting": userSetting,
	}).Debug("UserSetting found")

	c.JSON(http.StatusOK, userSetting)
}

// updateUserSetting godoc
// @Summary Updates user userSetting
// @Accept  json
// @Produce  json
// @Param request body model.UserSetting true "UserSetting"
// @Success 200 {object} model.UserSetting
// @Security ApiKeyAuth
// @Router /private/user/setting [put]
func updateUserSetting(c *gin.Context) {

	var userSetting *model.UserSetting
	err := c.BindJSON(&userSetting)
	if err != nil {
		log.Info("Invalid body content on updateUserSetting")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"userSetting": userSetting,
	}).Info("Updating userSetting")

	userSetting = ds.UpdateUserSetting(userSetting)

	if userSetting == nil {
		log.Info("Invalid User setting when saving, missing ID?")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"userSetting": userSetting,
	}).Info("Updated userSetting")

	c.JSON(http.StatusOK, userSetting)
}

// createGameDefinition godoc
// @Summary create Game definition
// @Accept json
// @Produce json
// @Param request body model.GameDefinition true "GameDefinition"
// @Success 200 {object} model.GameDefinition
// @Security ApiKeyAuth
// @Router /internal/game-definition [post]
func createGameDefinition(c *gin.Context) {

	var gameDefinition *model.GameDefinition
	err := c.BindJSON(&gameDefinition)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Invalid body content on createGameDefinition")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"gameDefinition": gameDefinition,
	}).Info("createGameDefinition")

	createResult := ds.CreateGameDefinition(gameDefinition)
	if createResult == nil {
		log.Error("Invalid GameDefinition when saving")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// load all the fields
	result := ds.FindGameDefinition(createResult.ID)

	log.WithFields(log.Fields{
		"gameDefinition": result,
		"ID":             createResult.ID,
	}).Info("createGameDefinition after create")

	c.JSON(http.StatusOK, result)
}

// updateGameDefinition godoc
// @Summary update Game definition
// @Accept json
// @Produce json
// @Param request body model.GameDefinition true "GameDefinition"
// @Success 200 {object} model.GameDefinition
// @Security ApiKeyAuth
// @Router /internal/game-definition [put]
func updateGameDefinition(c *gin.Context) {

	var gameDefinition *model.GameDefinition
	err := c.BindJSON(&gameDefinition)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Invalid body content on createGameDefinition")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"gameDefinition": gameDefinition,
	}).Info("updateGameDefinition")

	result := ds.UpdateGameDefinition(gameDefinition)
	if result == nil {
		log.Error("Invalid GameDefinition when updating")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, result)
}

// getUser godoc
// @Summary find The current user information
// @Accept json
// @Produce json
// @Success 200 {object} model.UserDetails
// @Security ApiKeyAuth
// @Router /private/get-user [get]
func getUser(c *gin.Context) {
	result := httphelper.UserDetailsFromContext(c)
	result.Classrooms = ds.FindAllClassroomByStudent(result.User.ID)
	result.Settings = *ds.FindUserSettingByUser(result.User)
	result.Level = *ds.FindUserLevelByUserID(result.User.ID)

	log.WithFields(log.Fields{
		"user": result,
	}).Error("get-user")

	c.JSON(http.StatusOK, result)
}

// getUserDashboard godoc
// @Summary find The current user information
// @Accept json
// @Produce json
// @Success 200 {object} model.UserDetails
// @Security ApiKeyAuth
// @Router /dashboard/get-user [get]
func getUserDashboard(c *gin.Context) {
	result := httphelper.UserDetailsFromContext(c)
	c.JSON(http.StatusOK, result)
}

// getLuchador godoc
// @Summary find or create Luchador for the current user
// @Accept json
// @Produce json
// @Success 200 {object} model.GameComponent
// @Security ApiKeyAuth
// @Router /private/luchador [get]
func getLuchador(c *gin.Context) {
	user := httphelper.UserFromContext(c)
	var luchador *model.GameComponent

	luchador = ds.FindLuchador(user)
	log.WithFields(log.Fields{
		"luchador": model.LogGameComponent(luchador),
		"user.id":  user.ID,
	}).Info("after find luchador on getLuchador")

	if luchador == nil {
		luchador = &model.GameComponent{
			UserID: user.ID,
			Name:   fmt.Sprintf("Luchador%d", user.ID),
		}

		luchador.Configs = model.RandomConfig()
		luchador.Name = model.RandomName(luchador.Configs)
		log.WithFields(log.Fields{
			"getLuchador": model.LogGameComponent(luchador),
		}).Info("creating luchador")

		luchador = ds.CreateLuchador(luchador)

		if luchador == nil {
			log.Error("Invalid Luchador when saving")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		log.WithFields(log.Fields{
			"luchador": model.LogGameComponent(luchador),
		}).Info("created luchador")
	}

	log.WithFields(log.Fields{
		"getLuchador": model.LogGameComponent(luchador),
	}).Info("result")

	c.JSON(http.StatusOK, luchador)
}

// updateLuchador godoc
// @Summary Updates Luchador
// @Accept  json
// @Produce  json
// @Param request body model.GameComponent true "Luchador"
// @Success 200 {object} model.UpdateLuchadorResponse
// @Security ApiKeyAuth
// @Router /private/luchador [put]
func updateLuchador(c *gin.Context) {
	user := httphelper.UserFromContext(c)
	response := model.UpdateLuchadorResponse{Errors: []string{}}

	var luchador *model.GameComponent
	err := c.BindJSON(&luchador)
	if err != nil {
		log.Info("Invalid body content on updateLuchador")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	luchador.Name = cleanName(luchador.Name)

	if len(luchador.Name) < 3 {
		response.Errors = append(response.Errors, "Luchador name length should be at least 3 characters")
	}

	if len(luchador.Name) > 40 {
		response.Errors = append(response.Errors, "Luchador name length should be less or equal to 40 characters")
	}

	if utility.ContainsBadWord(luchador.Name) {
		response.Errors = append(response.Errors, "Luchador name contains inappropriate language")
	}

	if ds.NameExist(luchador.ID, luchador.Name) {
		response.Errors = append(response.Errors, "Luchador with this name already exists")
	}

	if len(response.Errors) > 0 {
		response.Luchador = luchador

		log.WithFields(log.Fields{
			"luchador": luchador,
			"response": response,
		}).Debug("updateLuchador")

		c.JSON(http.StatusOK, response)
		return
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
		"action":   "before save",
	}).Debug("updateLuchador")

	// validate if the luchador is the same from the user
	currentLuchador := ds.FindLuchador(user)
	log.WithFields(log.Fields{
		"luchador": luchador,
		"user.ID":  user.ID,
	}).Info("find luchador for current user")

	if luchador.ID != currentLuchador.ID {
		log.Info("Invalid Luchador.ID on updateLuchador")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response.Luchador = ds.UpdateLuchador(luchador)

	if response.Luchador == nil {
		log.Info("Invalid Luchador when saving, missing ID?")
		response.Errors = append(response.Errors, "Invalid Luchador when saving, missing ID?")
	} else {
		channel := fmt.Sprintf("luchador.%v.update", response.Luchador.ID)
		luchadorUpdateJSON, _ := json.Marshal(response.Luchador)
		message := string(luchadorUpdateJSON)
		publisher.Publish(channel, message)
	}

	log.WithFields(log.Fields{
		"response": response,
		"action":   "after save",
		"errors":   response.Errors,
	}).Info("updateLuchador")

	c.JSON(http.StatusOK, response)
}

func cleanName(name string) string {
	name = strings.TrimSpace(name)
	return name
}

// getTutorialGameDefinition godoc
// @Summary find tutorial GameDefinition
// @Accept json
// @Produce json
// @Success 200 {array} model.GameDefinition
// @Security ApiKeyAuth
// @Router /private/tutorial [get]
func getTutorialGameDefinition(c *gin.Context) {

	tutorials := ds.FindTutorialGameDefinition()

	log.WithFields(log.Fields{
		"tutorials": tutorials,
	}).Info("getTutorialGameDefinition")

	c.JSON(http.StatusOK, tutorials)
}

// getMaskConfig godoc
// @Summary find maskConfig for a luchador
// @Accept json
// @Produce json
// @Param id path int true "Luchador ID"
// @Success 200 {array} model.Config
// @Security ApiKeyAuth
// @Router /private/mask-config/{id} [get]
func getMaskConfig(c *gin.Context) {

	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Info("Invalid ID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"id": aid,
	}).Info("getMaskConfig")

	configs := ds.FindMaskConfig(uint(aid))

	log.WithFields(log.Fields{
		"configs": configs,
	}).Info("getMaskConfig")

	c.JSON(http.StatusOK, configs)
}

// getGameDefinitionByName godoc
// @Summary find a game definition
// @Accept json
// @Produce json
// @Param name path string true "GameDefinition name"
// @Success 200 {object} model.GameDefinition
// @Security ApiKeyAuth
// @Router /internal/game-definition/{name} [get]
func getGameDefinitionByName(c *gin.Context) {

	name := c.Param("name")

	log.WithFields(log.Fields{
		"name": name,
	}).Info("getGameDefinition")

	gameDefinition := ds.FindGameDefinitionByName(name)

	log.WithFields(log.Fields{
		"gameDefinition": gameDefinition,
	}).Info("getGameDefinition")

	c.JSON(http.StatusOK, gameDefinition)
}

// getGameDefinitionByIDInternal godoc
// @Summary find a game definition
// @Accept json
// @Produce json
// @Param id path int true "GameDefinition id"
// @Success 200 {object} model.GameDefinition
// @Security ApiKeyAuth
// @Router /internal/game-definition-id/{id} [get]
func getGameDefinitionByIDInternal(c *gin.Context) {
	getGameDefinitionByID(c)
}

// getGameDefinitionByID godoc
// @Summary find a game definition
// @Accept json
// @Produce json
// @Param id path int true "GameDefinition id"
// @Success 200 {object} model.GameDefinition
// @Security ApiKeyAuth
// @Router /private/game-definition-id/{id} [get]
func getGameDefinitionByID(c *gin.Context) {

	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Info("Invalid ID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"id": aid,
	}).Info("getGameDefinitionByID")

	gameDefinition := ds.FindGameDefinition(uint(aid))

	log.WithFields(log.Fields{
		"gameDefinition": gameDefinition,
	}).Info("getGameDefinition")

	c.JSON(http.StatusOK, gameDefinition)
}

// getGameDefinition godoc
// @Summary find all game definitions
// @Accept json
// @Produce json
// @Success 200 {array} model.GameDefinition
// @Security ApiKeyAuth
// @Router /private/game-definition-all [get]
func getGameDefinition(c *gin.Context) {

	result := ds.FindAllSystemGameDefinition()

	log.WithFields(log.Fields{
		"result": result,
	}).Info("getGameDefinition")

	c.JSON(http.StatusOK, result)
}

// getRandomMaskConfig godoc
// @Summary create random maskConfig
// @Accept json
// @Produce json
// @Success 200 {array} model.Config
// @Security ApiKeyAuth
// @Router /private/mask-random [get]
func getRandomMaskConfig(c *gin.Context) {

	log.Info("getRandomMaskConfig")
	configs := model.RandomConfig()

	log.WithFields(log.Fields{
		"configs": configs,
	}).Info("getRandomMaskConfig")

	c.JSON(http.StatusOK, configs)
}

// createGameComponent godoc
// @Summary Create Gamecomponent as Luchador
// @Accept  json
// @Produce  json
// @Param request body model.GameComponent true "Luchador"
// @Success 200 {object} model.GameComponent
// @Security ApiKeyAuth
// @Router /internal/game-component [post]
func createGameComponent(c *gin.Context) {

	var luchador *model.GameComponent
	err := c.BindJSON(&luchador)
	if err != nil {
		log.Info("Invalid body content on createGameComponent")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
		"action":   "before save",
	}).Info("createGameComponent")

	// validate if the luchador is the same from the user

	found := ds.FindLuchadorByName(luchador.Name)

	if found == nil {
		log.Info("Luchador not found, will create")
		luchador.Configs = model.RandomConfig()
		log.WithFields(log.Fields{
			"configs": luchador.Configs,
		}).Info("Random config assigned to luchador")

		luchador = ds.CreateLuchador(luchador)
		luchador = ds.FindLuchadorByID(luchador.ID)
	} else {
		luchador = found
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
		"action":   "after save",
	}).Info("createGameComponent")

	c.JSON(http.StatusOK, luchador)
}

// getActiveMatches godoc
// @Summary find active matches
// @Accept json
// @Produce json
// @Success 200 {array} model.ActiveMatch
// @Security ApiKeyAuth
// @Router /private/match [get]
func getActiveMatches(c *gin.Context) {

	var result []model.ActiveMatch

	// multiplayer matches
	matches := *ds.FindActiveMultiplayerMatches()
	for _, match := range matches {
		gameDefinition := ds.FindGameDefinition(match.GameDefinitionID)
		add := model.ActiveMatch{
			MatchID:     match.ID,
			Name:        gameDefinition.Name,
			Label:       gameDefinition.Label,
			Description: gameDefinition.Description,
			Type:        gameDefinition.Type,
			SortOrder:   gameDefinition.SortOrder,
			Duration:    gameDefinition.Duration,
			TimeStart:   match.TimeStart,
		}

		result = append(result, add)
	}

	// gamedefinitions
	gameDefinitions := *ds.FindTutorialGameDefinition()
	for _, gameDefinition := range gameDefinitions {
		add := model.ActiveMatch{
			MatchID:     0,
			Name:        gameDefinition.Name,
			Label:       gameDefinition.Label,
			Description: gameDefinition.Description,
			Type:        gameDefinition.Type,
			SortOrder:   gameDefinition.SortOrder,
			Duration:    gameDefinition.Duration,
		}

		result = append(result, add)
	}

	log.WithFields(log.Fields{
		"matches": result,
	}).Info("getActiveMatches")

	c.JSON(http.StatusOK, &result)
}

// getActiveMultiplayerMatches godoc
// @Summary find active multiplayer matches
// @Accept json
// @Produce json
// @Success 200 {array} model.Match
// @Security ApiKeyAuth
// @Router /private/match-multiplayer [get]
func getActiveMultiplayerMatches(c *gin.Context) {

	// multiplayer matches
	matches := ds.FindActiveMultiplayerMatches()

	log.WithFields(log.Fields{
		"matches": model.LogMatches(matches),
	}).Info("getActiveMultiplayerMatches")

	c.JSON(http.StatusOK, matches)
}

// getMatchInternal godoc
// @Summary find one match
// @Accept json
// @Produce json
// @Param matchID query int false "int valid"
// @Success 200 {object} model.Match
// @Security ApiKeyAuth
// @Router /internal/match-single [get]
func getMatchInternal(c *gin.Context) {
	getMatch(c)
}

// getMatch godoc
// @Summary find one match
// @Accept json
// @Produce json
// @Param matchID query int false "int valid"
// @Success 200 {object} model.Match
// @Security ApiKeyAuth
// @Router /private/match-single [get]
func getMatch(c *gin.Context) {
	parameter := c.Query("matchID")
	i32, err := strconv.ParseInt(parameter, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"matchID": parameter,
		}).Error("Invalid matchID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var matchID uint
	matchID = uint(i32)

	log.WithFields(log.Fields{
		"matchID": matchID,
	}).Info("getMatch")

	match := ds.FindMatch(matchID)

	log.WithFields(log.Fields{
		"match": match,
	}).Info("getMatch")

	c.JSON(http.StatusOK, match)
}

// getMatchScore godoc
// @Summary find one match score
// @Accept json
// @Produce json
// @Param matchID query int false "int valid"
// @Success 200 {array} model.MatchScore
// @Security ApiKeyAuth
// @Router /private/match-score [get]
func getMatchScore(c *gin.Context) {

	parameter := c.Query("matchID")
	i32, err := strconv.ParseInt(parameter, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"matchID": parameter,
		}).Error("Invalid matchID on getMatchScore")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var matchID uint
	matchID = uint(i32)

	log.WithFields(log.Fields{
		"matchID": matchID,
	}).Info("getMatchScore")

	scores := ds.GetMatchScoresByMatchID(matchID)

	log.WithFields(log.Fields{
		"scores": scores,
	}).Info("getMatchScores")

	c.JSON(http.StatusOK, scores)
}

// getLuchadorConfigsForCurrentMatch godoc
// @Summary return luchador configs for current match
// @Accept json
// @Produce json
// @Param matchID query int false "int valid"
// @Success 200 {array} model.GameComponent
// @Security ApiKeyAuth
// @Router /private/match-config [get]
func getLuchadorConfigsForCurrentMatch(c *gin.Context) {

	parameter := c.Query("matchID")
	i32, err := strconv.ParseInt(parameter, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"matchID": parameter,
		}).Error("Invalid matchID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var matchID uint
	matchID = uint(i32)

	var result *[]model.GameComponent

	result = ds.FindLuchadorConfigsByMatchID(matchID)
	log.WithFields(log.Fields{
		"result": result,
	}).Debug("getLuchadorConfigsForCurrentMatch")

	c.JSON(http.StatusOK, result)
}

// joinMatch godoc
// @Summary Sends message with the request to join the match
// @Accept json
// @Produce json
// @Param request body model.JoinMatch true "JoinMatch"
// @Success 200 {object} model.Match
// @Security ApiKeyAuth
// @Router /private/join-match [post]
func joinMatch(c *gin.Context) {

	var joinMatch *model.JoinMatch
	err := c.BindJSON(&joinMatch)
	if err != nil {
		log.Info("Invalid body content on joinMatch")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := httphelper.UserFromContext(c)

	var luchador *model.GameComponent
	luchador = ds.FindLuchador(user)
	if luchador == nil {
		log.WithFields(log.Fields{
			"user": user,
		}).Error("Error getting luchador for the current user")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// make sure it will join with the luchador associated with the user
	joinMatch.LuchadorID = luchador.ID

	var match *model.Match
	match = ds.FindMatch(joinMatch.MatchID)
	if match == nil {
		log.WithFields(log.Fields{
			"user": user,
		}).Error("Match not found when trying to join match")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
		"match":    match,
	}).Info("Joining match")

	channel := fmt.Sprintf("match.%v.join", joinMatch.MatchID)
	joinMatchJSON, _ := json.Marshal(joinMatch)
	message := string(joinMatchJSON)

	publisher.Publish(channel, message)

	c.JSON(http.StatusOK, match)
}

// getLuchadorByID godoc
// @Summary find Luchador by ID
// @Accept json
// @Produce json
// @Param request body model.FindLuchadorWithGamedefinition true "FindLuchadorWithGamedefinition"
// @Success 200 {object} model.GameComponent
// @Security ApiKeyAuth
// @Router /internal/luchador [post]
func getLuchadorByIDAndGamedefinitionID(c *gin.Context) {

	var parameters *model.FindLuchadorWithGamedefinition
	err := c.BindJSON(&parameters)
	if err != nil {
		log.Info("Invalid body content on getLuchadorByIDAndGamedefinitionID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var luchador *model.GameComponent
	luchador = ds.FindLuchadorByID(parameters.LuchadorID)

	if luchador == nil {
		log.WithFields(log.Fields{
			"luchadorID": parameters.LuchadorID,
		}).Error("Luchador not found")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	filteredCodes := make([]model.Code, 0)
	for _, code := range luchador.Codes {
		if code.GameDefinitionID == parameters.GameDefinitionID {
			filteredCodes = append(filteredCodes, code)
		}
	}
	luchador.Codes = filteredCodes

	log.WithFields(log.Fields{
		"getLuchador": model.LogGameComponent(luchador),
	}).Info("result")

	c.JSON(http.StatusOK, luchador)
}

// getHealth godoc
// @Summary returns application health check information
// @Success 200
// @Security ApiKeyAuth
// @Router /internal/ready [get]
func getReady(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// addMatchPartipant godoc
// @Summary Adds luchador to a match
// @Accept json
// @Produce json
// @Param request body model.MatchParticipant true "MatchParticipant"
// @Success 200 {object} model.MatchParticipant
// @Security ApiKeyAuth
// @Router /internal/match-participant [post]
func addMatchPartipant(c *gin.Context) {

	var matchParticipantRequest *model.MatchParticipant
	err := c.BindJSON(&matchParticipantRequest)
	if err != nil {
		log.Info("Invalid body content on addMatchPartipant")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	matchParticipant := ds.AddMatchParticipant(matchParticipantRequest)
	if matchParticipant == nil {
		log.WithFields(log.Fields{
			"matchParticipant": matchParticipantRequest,
		}).Error("Error saving matchParticipant")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"matchParticipant": matchParticipant,
	}).Info("result")

	c.JSON(http.StatusOK, matchParticipant)
}

// endMatch godoc
// @Summary ends existing match
// @Accept json
// @Produce json
// @Param request body model.Match true "Match"
// @Success 200 {object} model.Match
// @Security ApiKeyAuth
// @Router /internal/end-match [put]
func endMatch(c *gin.Context) {

	var matchRequest *model.Match
	err := c.BindJSON(&matchRequest)
	if err != nil {
		log.Info("Invalid body content on endMatch")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	match := ds.EndMatch(matchRequest)
	if match == nil {
		log.WithFields(log.Fields{
			"match": matchRequest,
		}).Error("Error calling endMatch")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"match": match,
	}).Info("result")

	ds.UpdateParticipantsLevel(matchRequest.ID)

	c.JSON(http.StatusOK, match)
}

// runMatch godoc
// @Summary notify that the match is running, all participants joined
// @Accept json
// @Produce json
// @Param request body model.Match true "Match"
// @Success 200 {object} model.Match
// @Security ApiKeyAuth
// @Router /internal/run-match [put]
func runMatch(c *gin.Context) {

	var matchRequest *model.Match
	err := c.BindJSON(&matchRequest)
	if err != nil {
		log.Info("Invalid body content on runMatch")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	match := ds.RunMatch(matchRequest)
	if match == nil {
		log.WithFields(log.Fields{
			"match": matchRequest,
		}).Error("Error calling runMatch")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"match": match,
	}).Info("result")

	c.JSON(http.StatusOK, match)
}

// addMatchScore godoc
// @Summary saves a match score
// @Accept json
// @Produce json
// @Param request body model.ScoreList true "ScoreList"
// @Success 200 {object} model.MatchScore
// @Security ApiKeyAuth
// @Router /internal/add-match-scores [post]
func addMatchScores(c *gin.Context) {
	var scoreRequest *model.ScoreList
	err := c.BindJSON(&scoreRequest)
	if err != nil {
		log.Info("Invalid body content on addMatchScore")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	score := ds.AddMatchScores(scoreRequest)
	if score == nil {
		log.WithFields(log.Fields{
			"score": scoreRequest,
		}).Error("Error saving score")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"score": score,
	}).Info("result")

	c.JSON(http.StatusOK, score)
}

// addMatchMetric godoc
// @Summary saves a match metric
// @Accept json
// @Produce json
// @Param request body model.MatchMetric true "MatchMetric"
// @Success 200 {string} string
// @Security ApiKeyAuth
// @Router /internal/match-metric [post]
func addMatchMetric(c *gin.Context) {
	var metric *model.MatchMetric
	err := c.BindJSON(&metric)
	if err != nil {
		log.Info("Invalid body content on addMatchMetric")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result := eventsDS.AddMatchMetric(metric)
	if result == nil {
		log.WithFields(log.Fields{
			"metric": metric,
		}).Error("Error saving metric")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"metric": result,
	}).Debug("result")

	c.JSON(http.StatusOK, "")
}

// getClassroom godoc
// @Summary find all Classroom
// @Accept json
// @Produce json
// @Success 200 {array} model.Classroom
// @Security ApiKeyAuth
// @Router /dashboard/classroom [get]
func getClassroom(c *gin.Context) {
	user := httphelper.UserFromContext(c)
	result := ds.FindAllClassroom(user)

	log.WithFields(log.Fields{
		"result": result,
	}).Info("getClassroom")

	c.JSON(http.StatusOK, result)
}

// getClassroomStudents godoc
// @Summary find all Classroom students
// @Accept json
// @Produce json
// @Param id path int true "Classroom id"
// @Success 200 {array} model.StudentResponse
// @Security ApiKeyAuth
// @Router /dashboard/classroom/students/{id} [get]
func getClassroomStudents(c *gin.Context) {

	// get classroom id parameter
	id, err := httphelper.GetIntegerParam(c, "id", "getClassroomStudents")
	if err != nil {
		log.Info("Invalid body content on getClassroomStudents")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check if current user can see the list of students
	user := httphelper.UserFromContext(c)
	classroom := ds.FindClassroomByID(id)
	if classroom.OwnerID != user.ID {
		log.WithFields(log.Fields{
			"classroom": classroom,
			"user":      user,
		}).Warn("Current user is not the owner of this classroom")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// get the list of students
	// add user name, DONT add user email and name
	result := ds.BuildStudentResponse(classroom.Students)

	log.WithFields(log.Fields{
		"result": result,
	}).Info("getClassroom")

	c.JSON(http.StatusOK, result)
}

// addClassroom godoc
// @Summary add a Classroom
// @Accept json
// @Produce json
// @Param request body model.Classroom true "Classroom"
// @Success 200 {object} model.Classroom
// @Security ApiKeyAuth
// @Router /dashboard/classroom [post]
func addClassroom(c *gin.Context) {
	user := httphelper.UserFromContext(c)

	var classroom *model.Classroom
	err := c.BindJSON(&classroom)
	if err != nil {
		log.Info("Invalid body content on addClassroom")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	classroom.OwnerID = user.ID

	result := ds.AddClassroom(classroom)
	if result == nil {
		log.WithFields(log.Fields{
			"classroom": classroom,
		}).Error("Error saving classroom")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"classroom": result,
	}).Debug("results")

	c.JSON(http.StatusOK, result)
}

// joinClassroom godoc
// @Summary join a classroom
// @Accept json
// @Produce json
// @Param accessCode path string true "classroom access code"
// @Success 200 {object} model.Classroom
// @Security ApiKeyAuth
// @Router /private/join-classroom/{accessCode} [post]
func joinClassroom(c *gin.Context) {

	accessCode := c.Param("accessCode")
	user := httphelper.UserFromContext(c)

	log.WithFields(log.Fields{
		"accessCode": accessCode,
	}).Info("joinClassroom")

	classroom := ds.JoinClassroom(user, accessCode)

	log.WithFields(log.Fields{
		"classroom": classroom,
	}).Info("classroom")

	c.JSON(http.StatusOK, classroom)
}

// getPublicAvailableMatch godoc
// @Summary find all public available matches
// @Accept json
// @Produce json
// @Success 200 {array} model.AvailableMatch
// @Security ApiKeyAuth
// @Router /private/available-match-public [get]
func getPublicAvailableMatch(c *gin.Context) {
	result := ds.FindPublicAvailableMatch()

	log.WithFields(log.Fields{
		"result": model.LogAvailableMatches(result),
	}).Info("getPublicAvailableMatch")

	c.JSON(http.StatusOK, result)
}

// getClassroomAvailableMatch godoc
// @Summary find available matches by classroom
// @Accept json
// @Produce json
// @Param id path int true "Classroom id"
// @Success 200 {array} model.AvailableMatch
// @Security ApiKeyAuth
// @Router /private/available-match-classroom/{id} [get]
func getClassroomAvailableMatch(c *gin.Context) {
	id, err := httphelper.GetIntegerParam(c, "id", "getClassroomAvailableMatch")
	if err != nil {
		log.Info("Invalid body content on getClassroomAvailableMatch")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result := ds.FindAvailableMatchByClassroomID(id)

	log.WithFields(log.Fields{
		"result": model.LogAvailableMatches(result),
	}).Info("getPublicAvailableMatch")

	c.JSON(http.StatusOK, result)
}

// addEvents godoc
// @Summary add page events
// @Accept json
// @Produce json
// @Param request body model.PageEventRequest true "PageEventRequest"
// @Success 200 {string} string
// @Security ApiKeyAuth
// @Router /private/page-events [post]
func addEvents(c *gin.Context) {
	user := httphelper.UserFromContext(c)

	var request *model.PageEventRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Info("Invalid body content on addEvent")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	header := c.Request.Header.Get("User-Agent")
	ua := useragent.Parse(header)

	event := model.PageEvent{
		UserID:      user.ID,
		RemoteAddr:  c.ClientIP(),
		UserAgent:   ua.Name,
		Version:     ua.Version,
		OSName:      ua.OS,
		OSVersion:   ua.OSVersion,
		Mobile:      ua.Mobile,
		Tablet:      ua.Tablet,
		Desktop:     ua.Desktop,
		Device:      ua.Device,
		Page:        request.Page,
		Action:      request.Action,
		ComponentID: request.ComponentID,
		AppName:     request.AppName,
		AppVersion:  request.AppVersion,
		Value1:      request.Value1,
		Value2:      request.Value2,
		Value3:      request.Value3,
	}

	eventsDS.CreateEvent(event)
	c.JSON(http.StatusOK, "")
}

// getLevelGroup godoc
// @Summary find all level groups
// @Accept json
// @Produce json
// @Success 200 {array} model.LevelGroup
// @Security ApiKeyAuth
// @Router /private/level-group [get]
func getLevelGroup(c *gin.Context) {
	result := ds.FindLevelGroup()

	log.WithFields(log.Fields{
		"result": result,
	}).Info("getLevelGroup")

	c.JSON(http.StatusOK, result)
}
