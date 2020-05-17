package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
)

const (
	cookieName             = "kc-access"
	getkeeperEncryptionKey = "GATEKEEPER_ENCRYPTION_KEY"
	dashboardRole          = "dashboard_user"
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
		level := ds.FindUserLevelByUserID(user.ID)

		userDetails := model.UserDetails{
			User:  user,
			Roles: sessionUser.Roles,
			Level: *level,
		}
		c.Set("userDetails", userDetails)
	}
}

// SessionIsValidAndDashBoardUser check if Authorization header is valid
func SessionIsValidAndDashBoardUser(ds *datasource.DataSource) gin.HandlerFunc {
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

		if !contains(sessionUser.Roles, dashboardRole) {
			log.WithFields(log.Fields{
				"sessionUser":   sessionUser,
				"dashboardRole": dashboardRole,
			}).Info("User DONT have dashboard role")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		user := ds.CreateUser(sessionUser.Username)
		level := ds.FindUserLevelByUserID(user.ID)

		userDetails := model.UserDetails{
			User:  user,
			Roles: sessionUser.Roles,
			Level: *level,
		}
		c.Set("userDetails", userDetails)
	}
}

func contains(roles []string, search string) bool {
	for _, role := range roles {
		if role == search {
			return true
		}
	}
	return false
}

// SessionAllwaysValid test function for local development
func SessionAllwaysValid(ds *datasource.DataSource) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := ds.CreateUser("test")
		level := ds.FindUserLevelByUserID(user.ID)
		userDetails := model.UserDetails{
			User:  user,
			Roles: []string{dashboardRole},
			Level: *level,
		}
		c.Set("userDetails", userDetails)
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
