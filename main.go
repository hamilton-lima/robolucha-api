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
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gitlab.com/robolucha/robolucha-api/auth"

	_ "gitlab.com/robolucha/robolucha-api/docs"
)

//UpdateLuchadorResponse data structure
type UpdateLuchadorResponse struct {
	Errors   []string       `json:"errors"`
	Luchador *GameComponent `json:"luchador"`
}

var dataSource *DataSource
var publisher Publisher

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	rand.Seed(time.Now().UTC().UnixNano())

	log.Info("Robolucha API, start.")

	dataSource = NewDataSource(BuildMysqlConfig())
	defer dataSource.db.Close()

	publisher = &RedisPublisher{}
	go dataSource.KeepAlive()

	if len(os.Args) < 2 {
		log.Error("Missing gamedefinition folder parameter")
		os.Exit(2)
	}

	gameDefinitionFolder := os.Args[1]
	SetupGameDefinitionFromFolder(gameDefinitionFolder)

	port := os.Getenv("API_PORT")
	if len(port) == 0 {
		port = "5000"
	}

	log.WithFields(log.Fields{
		"port": port,
	}).Debug("Port configuration")

	internalAPIKey := os.Getenv("INTERNAL_API_KEY")
	logRequestBody := os.Getenv("GIM_LOG_REQUEST_BODY")
	router := createRouter(internalAPIKey, logRequestBody, SessionIsValid)
	router.Run(":" + port)

	log.WithFields(log.Fields{
		"port": port,
	}).Debug("Server is ready")
}

func createRouter(internalAPIKey string, logRequestBody string,
	factory SessionValidatorFactory) *gin.Engine {

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
	internalAPI.Use(KeyIsValid(internalAPIKey))
	{
		internalAPI.GET("/game-definition/:name", getGameDefinitionByName)
		internalAPI.GET("/game-definition-id/:id", getGameDefinitionByIDInternal)
		internalAPI.POST("/game-definition", createGameDefinition)
		internalAPI.PUT("/game-definition", updateGameDefinition)

		internalAPI.POST("/start-match/:name", startMatch)
		internalAPI.POST("/game-component", createGameComponent)
		internalAPI.POST("/luchador", getLuchadorByIDAndGamedefinitionID)
		internalAPI.POST("/match-participant", addMatchPartipant)
		internalAPI.PUT("/end-match", endMatch)
		internalAPI.GET("/ready", getReady)
		internalAPI.POST("/add-match-scores", addMatchScores)
		internalAPI.GET("/match-single", getMatchInternal)
	}

	privateAPI := router.Group("/private")
	privateAPI.Use(factory())
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
		privateAPI.GET("/match-single", getMatch)
		privateAPI.GET("/match-config", getLuchadorConfigsForCurrentMatch)
		privateAPI.POST("/join-match", joinMatch)
		privateAPI.GET("/game-definition-id/:id", getGameDefinitionByID)
		privateAPI.GET("/game-definition-all", getGameDefinition)
		privateAPI.POST("/start-tutorial-match/:name", startTutorialMatch)

	}

	return router
}

type SessionValidatorFactory func() gin.HandlerFunc

// SessionIsValid check if Authoraization header is valid
func SessionIsValid() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookieName := "kc-access"

		authorization, err := c.Request.Cookie(cookieName)
		if err != nil {
			log.Debug("Error reading authorization cookie")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if authorization.Value == "" {
			log.Debug("No Authorization cookie")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		key := os.Getenv("GATEKEEPER_ENCRYPTION_KEY")
		sessionUser, err := auth.GetUser(authorization.Value, key)
		if err != nil {
			log.Debug("Error reading user from authorization cookie")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if sessionUser.Username == "" {
			log.WithFields(log.Fields{
				"authorization": authorization,
				"cookie-name":   cookieName,
				"sessionUser":   sessionUser,
			}).Info("Invalid Session")
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			log.WithFields(log.Fields{
				"sessionUser": sessionUser,
			}).Info("User Authorized")
		}

		user := dataSource.createUser(User{Username: sessionUser.Username})
		c.Set("user", user)
	}
}

// KeyIsValid check if Authoraization header is valid
func KeyIsValid(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			log.Debug("No Authorization header")
			c.AbortWithStatus(http.StatusForbidden)
		}

		if authorization != key {
			log.Info("INVALID Authorization key")
			c.AbortWithStatus(http.StatusForbidden)
		}

		log.Info("VALID Authorization key")
	}
}

// findUserSetting godoc
// @Summary find current user userSetting
// @Accept  json
// @Produce  json
// @Success 200 {object} main.UserSetting
// @Security ApiKeyAuth
// @Router /private/user/setting [get]
func findUserSetting(c *gin.Context) {

	log.Info("Finding userSetting")
	val, _ := c.Get("user")
	user := val.(*User)

	userSetting := dataSource.findUserSettingByUser(user)

	log.WithFields(log.Fields{
		"userSetting": userSetting,
	}).Debug("UserSetting found")

	c.JSON(http.StatusOK, userSetting)
}

// updateUserSetting godoc
// @Summary Updates user userSetting
// @Accept  json
// @Produce  json
// @Param request body main.UserSetting true "UserSetting"
// @Success 200 {object} main.UserSetting
// @Security ApiKeyAuth
// @Router /private/user/setting [put]
func updateUserSetting(c *gin.Context) {

	var userSetting *UserSetting
	err := c.BindJSON(&userSetting)
	if err != nil {
		log.Info("Invalid body content on updateUserSetting")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"userSetting": userSetting,
	}).Info("Updating userSetting")

	userSetting = dataSource.updateUserSetting(userSetting)

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
// @Param request body main.GameDefinition true "GameDefinition"
// @Success 200 {object} main.GameDefinition
// @Security ApiKeyAuth
// @Router /internal/game-definition [post]
func createGameDefinition(c *gin.Context) {

	var gameDefinition *GameDefinition
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

	createResult := dataSource.createGameDefinition(gameDefinition)
	if createResult == nil {
		log.Error("Invalid GameDefinition when saving")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// load all the fields
	result := dataSource.findGameDefinition(createResult.ID)

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
// @Param request body main.GameDefinition true "GameDefinition"
// @Success 200 {object} main.GameDefinition
// @Security ApiKeyAuth
// @Router /internal/game-definition [put]
func updateGameDefinition(c *gin.Context) {

	var gameDefinition *GameDefinition
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

	result := dataSource.updateGameDefinition(gameDefinition)
	if result == nil {
		log.Error("Invalid GameDefinition when updating")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, result)
}

// startMatch godoc
// @Summary create Match
// @Accept json
// @Produce json
// @Param name path string true "GameDefinition name"
// @Success 200 {object} main.Match
// @Security ApiKeyAuth
// @Router /internal/start-match/{name} [post]
func startMatch(c *gin.Context) {

	name := c.Param("name")

	log.WithFields(log.Fields{
		"name": name,
	}).Info("startMatch")

	gameDefinition := dataSource.findGameDefinitionByName(name)
	if gameDefinition == nil {
		log.Info("Invalid gamedefinition name")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	match := dataSource.createMatch(gameDefinition.ID)
	if match == nil {
		log.Error("Invalid Match when saving")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// load all the fields
	match = dataSource.findMatch(match.ID)

	log.WithFields(log.Fields{
		"createMatch": match,
	}).Info("created match")

	c.JSON(http.StatusOK, match)
}

// startTutorialMatch godoc
// @Summary create Match and publish
// @Accept json
// @Produce json
// @Param name path string true "GameDefinition name"
// @Success 200 {object} main.JoinMatch
// @Security ApiKeyAuth
// @Router /private/start-tutorial-match/{name} [post]
func startTutorialMatch(c *gin.Context) {

	name := c.Param("name")

	log.WithFields(log.Fields{
		"name": name,
	}).Info("startTutorialMatch")

	gameDefinition := dataSource.findGameDefinitionByName(name)
	if gameDefinition == nil {
		log.Info("Invalid gamedefinition name")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	val, _ := c.Get("user")
	user := val.(*User)

	luchador := dataSource.findLuchador(user)
	if luchador == nil {
		log.WithFields(log.Fields{
			"user": user,
		}).Error("Error getting luchador for the current user")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	match := dataSource.findActiveMatchesByGameDefinitionAndParticipant(gameDefinition, luchador)
	// not found will create
	if match == nil {
		match = dataSource.createMatch(gameDefinition.ID)
		if match == nil {
			log.Error("Invalid Match when saving")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	result := JoinMatch{MatchID: match.ID, LuchadorID: luchador.ID}
	// publish event to run the match
	resultJSON, _ := json.Marshal(result)
	publisher.Publish("start.match", string(resultJSON))

	log.WithFields(log.Fields{
		"createMatch": result,
	}).Info("created match")

	c.JSON(http.StatusOK, result)
}

// getUser godoc
// @Summary find The current user information
// @Accept json
// @Produce json
// @Success 200 {object} main.User
// @Security ApiKeyAuth
// @Router /private/get-user [get]
func getUser(c *gin.Context) {
	val, _ := c.Get("user")
	user := val.(*User)
	c.JSON(http.StatusOK, user)
}

// getLuchador godoc
// @Summary find or create Luchador for the current user
// @Accept json
// @Produce json
// @Success 200 {object} main.GameComponent
// @Security ApiKeyAuth
// @Router /private/luchador [get]
func getLuchador(c *gin.Context) {
	val, _ := c.Get("user")
	user := val.(*User)
	var luchador *GameComponent

	luchador = dataSource.findLuchador(user)
	log.WithFields(log.Fields{
		"luchador": luchador,
		"user.id":  user.ID,
	}).Info("after find luchador on getLuchador")

	if luchador == nil {
		luchador = &GameComponent{
			UserID: user.ID,
			Name:   fmt.Sprintf("Luchador%d", user.ID),
		}

		luchador.Configs = randomConfig()
		luchador.Name = randomName(luchador.Configs)
		log.WithFields(log.Fields{
			"getLuchador": luchador,
		}).Info("creating luchador")

		luchador = dataSource.createLuchador(luchador)

		if luchador == nil {
			log.Error("Invalid Luchador when saving")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		log.WithFields(log.Fields{
			"luchador": luchador,
		}).Info("created luchador")
	}

	log.WithFields(log.Fields{
		"getLuchador": luchador,
	}).Info("result")

	c.JSON(http.StatusOK, luchador)
}

// updateLuchador godoc
// @Summary Updates Luchador
// @Accept  json
// @Produce  json
// @Param request body main.GameComponent true "Luchador"
// @Success 200 {object} main.UpdateLuchadorResponse
// @Security ApiKeyAuth
// @Router /private/luchador [put]
func updateLuchador(c *gin.Context) {
	val, _ := c.Get("user")
	user := val.(*User)
	response := UpdateLuchadorResponse{Errors: []string{}}

	var luchador *GameComponent
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

	if dataSource.NameExist(luchador.ID, luchador.Name) {
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
	currentLuchador := dataSource.findLuchador(user)
	log.WithFields(log.Fields{
		"luchador": luchador,
		"user.ID":  user.ID,
	}).Info("find luchador for current user")

	if luchador.ID != currentLuchador.ID {
		log.Info("Invalid Luchador.ID on updateLuchador")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response.Luchador = dataSource.updateLuchador(luchador)

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
// @Success 200 200 {array} main.GameDefinition
// @Security ApiKeyAuth
// @Router /private/tutorial [get]
func getTutorialGameDefinition(c *gin.Context) {

	tutorials := dataSource.findTutorialGameDefinition()

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
// @Success 200 200 {array} main.Config
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

	configs := dataSource.findMaskConfig(uint(aid))

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
// @Success 200 200 {object} main.GameDefinition
// @Security ApiKeyAuth
// @Router /internal/game-definition/{name} [get]
func getGameDefinitionByName(c *gin.Context) {

	name := c.Param("name")

	log.WithFields(log.Fields{
		"name": name,
	}).Info("getGameDefinition")

	gameDefinition := dataSource.findGameDefinitionByName(name)

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
// @Success 200 200 {object} main.GameDefinition
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
// @Success 200 200 {object} main.GameDefinition
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

	gameDefinition := dataSource.findGameDefinition(uint(aid))

	log.WithFields(log.Fields{
		"gameDefinition": gameDefinition,
	}).Info("getGameDefinition")

	c.JSON(http.StatusOK, gameDefinition)
}

// getGameDefinition godoc
// @Summary find all game definitions
// @Accept json
// @Produce json
// @Success 200 200 {array} main.GameDefinition
// @Security ApiKeyAuth
// @Router /private/game-definition-all [get]
func getGameDefinition(c *gin.Context) {

	result := dataSource.findAllGameDefinition()

	log.WithFields(log.Fields{
		"result": result,
	}).Info("getGameDefinition")

	c.JSON(http.StatusOK, result)
}

// getRandomMaskConfig godoc
// @Summary create random maskConfig
// @Accept json
// @Produce json
// @Success 200 200 {array} main.Config
// @Security ApiKeyAuth
// @Router /private/mask-random [get]
func getRandomMaskConfig(c *gin.Context) {

	log.Info("getRandomMaskConfig")
	configs := randomConfig()

	log.WithFields(log.Fields{
		"configs": configs,
	}).Info("getRandomMaskConfig")

	c.JSON(http.StatusOK, configs)
}

// createGameComponent godoc
// @Summary Create Gamecomponent as Luchador
// @Accept  json
// @Produce  json
// @Param request body main.GameComponent true "Luchador"
// @Success 200 {object} main.GameComponent
// @Security ApiKeyAuth
// @Router /internal/game-component [post]
func createGameComponent(c *gin.Context) {

	var luchador *GameComponent
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

	found := dataSource.findLuchadorByName(luchador.Name)

	if found == nil {
		log.Info("Luchador not found, will create")
		luchador.Configs = randomConfig()
		log.WithFields(log.Fields{
			"configs": luchador.Configs,
		}).Info("Random config assigned to luchador")

		luchador = dataSource.createLuchador(luchador)
		luchador = dataSource.findLuchadorByID(luchador.ID)
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
// @Success 200 {array} main.Match
// @Security ApiKeyAuth
// @Router /private/match [get]
func getActiveMatches(c *gin.Context) {

	var matches *[]Match

	matches = dataSource.findActiveMatches()
	log.WithFields(log.Fields{
		"matches": matches,
	}).Info("getActiveMatches")

	c.JSON(http.StatusOK, matches)
}

// getMatchInternal godoc
// @Summary find one match
// @Accept json
// @Produce json
// @Param matchID query int false "int valid"
// @Success 200 {object} main.Match
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
// @Success 200 {object} main.Match
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

	match := dataSource.findMatch(matchID)

	log.WithFields(log.Fields{
		"match": match,
	}).Info("getMatch")

	c.JSON(http.StatusOK, match)
}

// getLuchadorConfigsForCurrentMatch godoc
// @Summary return luchador configs for current match
// @Accept json
// @Produce json
// @Param matchID query int false "int valid"
// @Success 200 {array} main.GameComponent
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

	var result *[]GameComponent

	result = dataSource.findLuchadorConfigsByMatchID(matchID)
	log.WithFields(log.Fields{
		"result": result,
	}).Debug("getLuchadorConfigsForCurrentMatch")

	c.JSON(http.StatusOK, result)
}

// joinMatch godoc
// @Summary Sends message with the request to join the match
// @Accept json
// @Produce json
// @Param request body main.JoinMatch true "JoinMatch"
// @Success 200 {object} main.Match
// @Security ApiKeyAuth
// @Router /private/join-match [post]
func joinMatch(c *gin.Context) {

	var joinMatch *JoinMatch
	err := c.BindJSON(&joinMatch)
	if err != nil {
		log.Info("Invalid body content on joinMatch")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	val, _ := c.Get("user")
	user := val.(*User)

	var luchador *GameComponent
	luchador = dataSource.findLuchador(user)
	if luchador == nil {
		log.WithFields(log.Fields{
			"user": user,
		}).Error("Error getting luchador for the current user")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// make sure it will join with the luchador associated with the user
	joinMatch.LuchadorID = luchador.ID

	var match *Match
	match = dataSource.findMatch(joinMatch.MatchID)
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
// @Param request body main.FindLuchadorWithGamedefinition true "FindLuchadorWithGamedefinition"
// @Success 200 {object} main.GameComponent
// @Security ApiKeyAuth
// @Router /internal/luchador [post]
func getLuchadorByIDAndGamedefinitionID(c *gin.Context) {

	var parameters *FindLuchadorWithGamedefinition
	err := c.BindJSON(&parameters)
	if err != nil {
		log.Info("Invalid body content on getLuchadorByIDAndGamedefinitionID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var luchador *GameComponent
	luchador = dataSource.findLuchadorByID(parameters.LuchadorID)

	if luchador == nil {
		log.WithFields(log.Fields{
			"luchadorID": parameters.LuchadorID,
		}).Error("Luchador not found")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	filteredCodes := make([]Code, 0)
	for _, code := range luchador.Codes {
		if code.GameDefinitionID == parameters.GameDefinitionID {
			filteredCodes = append(filteredCodes, code)
		}
	}
	luchador.Codes = filteredCodes

	log.WithFields(log.Fields{
		"getLuchador": luchador,
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
// @Param request body main.MatchParticipant true "MatchParticipant"
// @Success 200 {object} main.MatchParticipant
// @Security ApiKeyAuth
// @Router /internal/match-participant [post]
func addMatchPartipant(c *gin.Context) {

	var matchParticipantRequest *MatchParticipant
	err := c.BindJSON(&matchParticipantRequest)
	if err != nil {
		log.Info("Invalid body content on addMatchPartipant")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	matchParticipant := dataSource.addMatchParticipant(matchParticipantRequest)
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
// @Param request body main.Match true "Match"
// @Success 200 {object} main.Match
// @Security ApiKeyAuth
// @Router /internal/end-match [put]
func endMatch(c *gin.Context) {

	var matchRequest *Match
	err := c.BindJSON(&matchRequest)
	if err != nil {
		log.Info("Invalid body content on endMatch")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	match := dataSource.endMatch(matchRequest)
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

	c.JSON(http.StatusOK, match)
}

// addMatchScore godoc
// @Summary saves a match score
// @Accept json
// @Produce json
// @Param request body main.ScoreList true "ScoreList"
// @Success 200 {object} main.MatchScore
// @Security ApiKeyAuth
// @Router /internal/add-match-scores [post]
func addMatchScores(c *gin.Context) {
	var scoreRequest *ScoreList
	err := c.BindJSON(&scoreRequest)
	if err != nil {
		log.Info("Invalid body content on addMatchScore")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	score := dataSource.addMatchScores(scoreRequest)
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
