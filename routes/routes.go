package routes

import (
	"github.com/gin-gonic/gin"
)

// RouterDefinition defines the common behavior of router definition
type RouterDefinition interface {
	Setup(group *gin.RouterGroup)
}

// Use setup the router definition to be used with the router group
func Use(group *gin.RouterGroup, definition RouterDefinition) {
	definition.Setup(group)
}
