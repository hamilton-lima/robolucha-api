package httphelper

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetIntegerParam(c *gin.Context, paramName string, context string) (uint, error) {

	parameter := c.Query(paramName)
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
