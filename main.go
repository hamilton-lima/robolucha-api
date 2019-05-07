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

	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "gitlab.com/robolucha/robolucha-api/docs"
)

// LoginRequest data structure
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse data structure
type LoginResponse struct {
	Error bool   `json:"error"`
	UUID  string `json:"uuid"`
}

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
	log.SetLevel(log.DebugLevel)

	log.Info("Robolucha API, start.")

	dataSource = NewDataSource(BuildMysqlConfig())
	defer dataSource.db.Close()

	publisher = &RedisPublisher{}

	AddTestUsers(dataSource)
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
	router := createRouter(internalAPIKey, logRequestBody)
	router.Run(":" + port)

	log.WithFields(log.Fields{
		"port": port,
	}).Debug("Server is ready")
}

func createRouter(internalAPIKey string, logRequestBody string) *gin.Engine {
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
		publicAPI.POST("/login", handleLogin)
		publicAPI.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	internalAPI := router.Group("/internal")
	internalAPI.Use(KeyIsValid(internalAPIKey))
	{
		internalAPI.POST("/match", createMatch)
		internalAPI.POST("/game-component", createGameComponent)
		internalAPI.GET("/luchador", getLuchadorByID)
		internalAPI.POST("/match-participant", addMatchPartipant)
		internalAPI.PUT("/end-match", endMatch)
		internalAPI.GET("/ready", getReady)
		internalAPI.POST("/add-match-scores", addMatchScores)
	}

	privateAPI := router.Group("/private")
	privateAPI.Use(SessionIsValid())
	{
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

// SessionIsValid check if Authoraization header is valid
func SessionIsValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			log.Debug("No Authorization header")
			c.AbortWithStatus(http.StatusForbidden)
		}

		user := dataSource.findUserBySession(authorization)
		if user == nil {
			log.WithFields(log.Fields{
				"UUID": authorization,
			}).Info("Invalid Session UUID")
			c.AbortWithStatus(http.StatusForbidden)
		}

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

// handleLogin godoc
// @Summary Logs the user
// @Accept  json
// @Produce  json
// @Param request body main.LoginRequest true "LoginRequest"
// @Success 200 {object} main.LoginResponse
// @Router /public/login [post]
func handleLogin(c *gin.Context) {

	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		log.Info("Invalid body content on Login")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"email": json.Email,
	}).Info("Login Attempt")

	response := LoginResponse{Error: true}
	user := dataSource.findUserByEmail(json.Email)
	log.WithFields(log.Fields{
		"user": user,
	}).Debug("User found after login")

	if user != nil {
		session := dataSource.createSession(user)
		response.Error = false
		response.UUID = session.UUID
	}

	c.JSON(http.StatusOK, response)
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
		"user":     user,
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

	if len(luchador.Name) < 3 {
		response.Errors = append(response.Errors, "Luchador name length should be at least 3 characters")
	}

	if len(luchador.Name) > 30 {
		response.Errors = append(response.Errors, "Luchador name length should be less or equal to 30 characters")
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
		"user":     user,
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
		"errors=========================================": response.Errors,
	}).Info("updateLuchador")

	c.JSON(http.StatusOK, response)
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
