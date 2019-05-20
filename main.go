// @title Robolucha API
// @version 1.0
// @description Robolucha API
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	Errors   []string  `json:"errors"`
	Luchador *Luchador `json:"luchador"`
}

var dataSource *DataSource
var publisher Publisher

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.Info("Robolucha API, start.")

	dataSource = NewDataSource(BuildMysqlConfig())
	defer dataSource.db.Close()

	publisher = &RedisPublisher{}
	go dataSource.KeepAlive()

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
		internalAPI.GET("/game-definition/:id", getGameDefinition)
		internalAPI.POST("/game-definition", createGameDefinition)
		internalAPI.POST("/match", createMatch)
		internalAPI.POST("/game-component", createGameComponent)
		internalAPI.GET("/luchador", getLuchadorByID)
		internalAPI.POST("/match-participant", addMatchPartipant)
		internalAPI.PUT("/end-match", endMatch)
		internalAPI.GET("/ready", getReady)
		internalAPI.POST("/add-match-scores", addMatchScores)
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
		privateAPI.GET("/match-config", getLuchadorConfigsForCurrentMatch)
		privateAPI.POST("/join-match", joinMatch)
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
			log.WithFields(log.Fields{
				"Authorization": authorization,
			}).Info("Invalid Authorization key")
			c.AbortWithStatus(http.StatusForbidden)
		}

		log.WithFields(log.Fields{
			"Authorization": authorization,
		}).Info(">> Authorization key")
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

// createMatch godoc
// @Summary create Match
// @Accept json
// @Produce json
// @Param request body main.Match true "Match"
// @Success 200 {object} main.Match
// @Security ApiKeyAuth
// @Router /internal/match [post]
func createMatch(c *gin.Context) {

	var match *Match
	err := c.BindJSON(&match)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Invalid body content on createMatch")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"createMatch": match,
	}).Info("creating match")

	match = dataSource.createMatch(match)
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
// @Success 200 {object} main.Luchador
// @Security ApiKeyAuth
// @Router /private/luchador [get]
func getLuchador(c *gin.Context) {
	val, _ := c.Get("user")
	user := val.(*User)
	var luchador *Luchador

	luchador = dataSource.findLuchador(user)
	log.WithFields(log.Fields{
		"luchador": luchador,
		"user.id":  user.ID,
	}).Info("after find luchador on getLuchador")

	if luchador == nil {
		luchador = &Luchador{
			UserID: user.ID,
			Name:   fmt.Sprintf("Luchador%d", user.ID),
		}

		luchador.Codes = defaultCode()
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
// @Param request body main.Luchador true "Luchador"
// @Success 200 {object} main.UpdateLuchadorResponse
// @Security ApiKeyAuth
// @Router /private/luchador [put]
func updateLuchador(c *gin.Context) {
	val, _ := c.Get("user")
	user := val.(*User)
	response := UpdateLuchadorResponse{Errors: []string{}}

	var luchador *Luchador
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

// getGameDefinition godoc
// @Summary find a game definition
// @Accept json
// @Produce json
// @Param name path string true "GameDefinition name"
// @Success 200 200 {array} main.GameDefinition
// @Security ApiKeyAuth
// @Router /internal/game-definition/{name} [get]
func getGameDefinition(c *gin.Context) {

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
// @Param request body main.Luchador true "Luchador"
// @Success 200 {object} main.Luchador
// @Security ApiKeyAuth
// @Router /internal/game-component [post]
func createGameComponent(c *gin.Context) {

	var luchador *Luchador
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

// getLuchadorConfigsForCurrentMatch godoc
// @Summary return luchador configs for current match
// @Accept json
// @Produce json
// @Param matchID query int false "int valid"
// @Success 200 {array} main.Luchador
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

	var result *[]Luchador

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

	var luchador *Luchador
	luchador = dataSource.findLuchador(user)
	if luchador == nil {
		log.WithFields(log.Fields{
			"user": user,
		}).Error("Error getting luchador for the current uses")
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
// @Param luchadorID query int false "int valid"
// @Success 200 {object} main.Luchador
// @Security ApiKeyAuth
// @Router /internal/luchador [get]
func getLuchadorByID(c *gin.Context) {

	parameter := c.Query("luchadorID")
	i32, err := strconv.ParseInt(parameter, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"luchadorID": parameter,
		}).Error("Invalid luchadorID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var luchadorID uint
	luchadorID = uint(i32)

	var luchador *Luchador

	luchador = dataSource.findLuchadorByID(luchadorID)
	if luchador == nil {
		log.WithFields(log.Fields{
			"luchadorID": luchadorID,
		}).Error("Luchador not found")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

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
