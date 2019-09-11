package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/datasource"
)

const (
	cookieName = "kc-access"
	getkeeperEncryptionKey = "GATEKEEPER_ENCRYPTION_KEY"
)

// SessionValidatorFactory definition
type SessionValidatorFactory func(ds *datasource.DataSource) gin.HandlerFunc

// SessionIsValid check if Authorization header is valid
func SessionIsValid(ds *datasource.DataSource) gin.HandlerFunc {
	return func(c *gin.Context) {

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

		key := os.Getenv(getkeeperEncryptionKey)
		sessionUser, err := GetUser(authorization.Value, key)
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

		user := ds.CreateUser(sessionUser.Username)
		c.Set("user", user)
	}
}

// SessionAllwaysValid test function for local development
func SessionAllwaysValid(ds *datasource.DataSource) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := ds.CreateUser("test")
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
