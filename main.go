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
	"net/http"
	"os"

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

var dataSource *DataSource

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// dataSource = NewDataSource(BuildMysqlConfig())
	// defer dataSource.db.Close()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	publicAPI := router.Group("/public")
	{
		publicAPI.POST("/login", handleLogin)
		publicAPI.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	internalAPI := router.Group("/internal")
	{
		internalAPI.POST("/match", createMatch)
	}

	privateAPI := router.Group("/private")
	privateAPI.Use(SessionIsValid())
	{
		privateAPI.PUT("/user/setting", updateUserSetting)
		privateAPI.GET("/user/setting", findUserSetting)
	}

	router.Run()

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
		log.Info("Invalid body content on createMatch")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"createMatch": match,
	}).Info("creating match")

	match = dataSource.createMatch(match)

	if match == nil {
		log.Info("Invalid Match when saving")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"createMatch": match,
	}).Info("created match")

	c.JSON(http.StatusOK, match)
}
