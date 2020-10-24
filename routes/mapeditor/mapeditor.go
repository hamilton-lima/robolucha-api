package mapeditor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/httphelper"
	"gitlab.com/robolucha/robolucha-api/pubsub"
)

// Init receive database and message queue objects
func Init(_ds *datasource.DataSource, _publisher pubsub.Publisher) *Router {
	requestHandler = NewRequestHandler(_ds, _publisher)

	return &Router{ds: _ds,
		publisher: _publisher,
	}
}

// RequestHandler definition
type RequestHandler struct {
	ds        *datasource.DataSource
	publisher pubsub.Publisher
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(_ds *datasource.DataSource, _publisher pubsub.Publisher) *RequestHandler {
	handler := RequestHandler{
		ds:        _ds,
		publisher: _publisher,
	}

	return &handler
}

var requestHandler *RequestHandler

// Router definition
type Router struct {
	ds        *datasource.DataSource
	publisher pubsub.Publisher
}

// Setup definition
func (router *Router) Setup(group *gin.RouterGroup) {
	group.GET("/mapeditor", getMyGameDefinitions)
	// group.POST("/mapeditor", addMyGameDefinition)
	// group.PATCH("/mapeditor/:id", updateMyGameDefinition)
	// group.DELETE("/mapeditor/:id", delMyGameDefinition)
}

// getMyGameDefinitions godoc
// @Summary find my gamedefitions
// @Accept json
// @Produce json
// @Success 200 {array} model.GameDefinition
// @Security ApiKeyAuth
// @Router /private/mapeditor [get]
func getMyGameDefinitions(c *gin.Context) {
	user := httphelper.UserDetailsFromContext(c)
	gameDefinitions := requestHandler.ds.FindGameDefinitionByOwner(user.User.ID)
	c.JSON(http.StatusOK, gameDefinitions)
}
