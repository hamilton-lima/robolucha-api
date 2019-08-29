package httphelper

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/model"
	"strconv"
)

// GetIntegerParam gets an integer parameter from request and validate
func GetIntegerParam(c *gin.Context, paramName string, context string) (uint, error) {

	parameter := c.Param(paramName)
	i32, err := strconv.ParseInt(parameter, 10, 32)
	if err != nil {
		message := "Invalid Integer parameter"
		log.WithFields(log.Fields{
			"parameterName": paramName,
			"parameter":     parameter,
			"context":       context,
		}).Error(message)
		return 0, errors.New(message)
	}

	result := uint(i32)
	return result, nil
}

// UserFromContext get the current user from the request context
func UserFromContext(c *gin.Context) *model.User {
	val, _ := c.Get("user")
	user := val.(*model.User)
	return user
}
