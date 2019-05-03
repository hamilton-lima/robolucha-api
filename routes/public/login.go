package public

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	_ "gitlab.com/robolucha/robolucha-api/context"
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

// handleLogin godoc
// @Summary Logs the user
// @Accept  json
// @Produce  json
// @Param request body main.LoginRequest true "LoginRequest"
// @Success 200 {object} main.LoginResponse
// @Router /public/login [post]
func HandleLogin(c *gin.Context) {

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
	user := Context.ds.findUserByEmail(json.Email)
	log.WithFields(log.Fields{
		"user": user,
	}).Debug("User found after login")

	if user != nil {
		session := Context.ds.createSession(user)
		response.Error = false
		response.UUID = session.UUID
	}

	c.JSON(http.StatusOK, response)
}
