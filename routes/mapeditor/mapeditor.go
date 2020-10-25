package mapeditor

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/httphelper"
	"gitlab.com/robolucha/robolucha-api/model"
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
	group.GET("/mapeditor/default", getDefaultGameDefinition)
	group.POST("/mapeditor", addMyGameDefinition)
	group.PUT("/mapeditor", updateMyGameDefinition)
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
	gameDefinitions := requestHandler.Find(user.User.ID)
	c.JSON(http.StatusOK, gameDefinitions)
}

// getDefaultGameDefinition godoc
// @ID getDefaultGameDefinition
// @Summary get default game definition
// @Accept json
// @Produce json
// @Success 200 {object} model.GameDefinition
// @Security ApiKeyAuth
// @Router /private/mapeditor/default [get]
func getDefaultGameDefinition(c *gin.Context) {
	defaultGameDefinition := requestHandler.GetDefault()
	c.JSON(http.StatusOK, defaultGameDefinition)
}

// addMyGameDefinition godoc
// @Summary add a single gamedefition for this user
// @Accept json
// @Produce json
// @Param request body model.GameDefinition true "GameDefinition"
// @Success 200 {string} string
// @Security ApiKeyAuth
// @Router /private/mapeditor [post]
func addMyGameDefinition(c *gin.Context) {
	user := httphelper.UserDetailsFromContext(c)

	// parse body parameter
	var gameDefinition *model.GameDefinition
	err := c.BindJSON(&gameDefinition)
	if err != nil {
		log.Info("Invalid body content on addMyGameDefinition")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = requestHandler.Add(user.User.ID, gameDefinition)
	if err != nil {
		c.AbortWithStatus(http.StatusConflict)
	} else {
		c.JSON(http.StatusOK, "")
	}

}

// updateMyGameDefinition godoc
// @Summary update gamedefition for this user
// @Accept json
// @Produce json
// @Param request body model.GameDefinition true "GameDefinition"
// @Success 200 {array} model.GameDefinition
// @Security ApiKeyAuth
// @Router /private/mapeditor [put]
func updateMyGameDefinition(c *gin.Context) {
	user := httphelper.UserDetailsFromContext(c)

	// parse body parameter
	var gameDefinition *model.GameDefinition
	err := c.BindJSON(&gameDefinition)
	if err != nil {
		log.Info("Invalid body content on addMyGameDefinition")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = requestHandler.Update(user.User.ID, gameDefinition)
	if err != nil {
		c.AbortWithStatus(http.StatusConflict)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

// Find godoc
func (handler *RequestHandler) Find(userID uint) *[]model.GameDefinition {
	return handler.ds.FindGameDefinitionByOwner(userID)
}

// GetDefault godoc
func (handler *RequestHandler) GetDefault() model.GameDefinition {
	return model.BuildDefaultGameDefinition()
}

// Add godoc
func (handler *RequestHandler) Add(userID uint, gameDefinition *model.GameDefinition) error {
	foundByName := handler.ds.FindGameDefinitionByName(gameDefinition.Name)
	if foundByName != nil {
		log.WithFields(log.Fields{
			"gameDefinition": foundByName,
		}).Info("gamedefinition already EXISTS with this name")
		return errors.New("Gamedefinition name already exists")
	} else {
		gameDefinition.ID = 0
		gameDefinition.OwnerUserID = userID

		createResult := handler.ds.CreateGameDefinition(gameDefinition)
		log.WithFields(log.Fields{
			"gameDefinition": createResult,
		}).Info("gamedefinition CREATED")
		return nil
	}
}

// Update godoc
func (handler *RequestHandler) Update(userID uint, gameDefinition *model.GameDefinition) error {
	foundByID := handler.ds.FindGameDefinition(gameDefinition.ID)
	// must exist to be updated
	if foundByID == nil {
		log.WithFields(log.Fields{}).Info("gamedefinition DOES NOT EXIST, cant be updated")
		return errors.New("Gamedefinition DOES NOT exist")
	}

	// must be the owner to update it
	if foundByID.OwnerUserID != userID {
		log.WithFields(log.Fields{
			"foundByID.OwnerUserID": foundByID.OwnerUserID,
			"userID":                userID,
		}).Info("current user dont OWNS this gamedefinition, cant be updated")
		return errors.New("current user DOES NOT OWN this Gamedefinition")
	}

	foundByName := handler.ds.FindGameDefinitionByName(gameDefinition.Name)
	// name must not exist
	if foundByName != nil && foundByName.ID != gameDefinition.ID {
		log.WithFields(log.Fields{
			"gameDefinition": foundByName,
		}).Info("gamedefinition already EXISTS with this name")
		return errors.New("gamedefinition already EXISTS with this name")

	} else {
		gameDefinition.OwnerUserID = userID

		updateResult := handler.ds.UpdateGameDefinition(gameDefinition)
		log.WithFields(log.Fields{
			"gameDefinition": updateResult,
		}).Info("gamedefinition UPDATED")
		return nil
	}

}
