package mapeditor

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gitlab.com/robolucha/robolucha-api/auth"
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
	group.POST("/mapeditor/update-classroom-map-availability", updateClassroomMapAvailability)
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

	if auth.UserBelongsToRole(user, auth.SystemEditorRole) {
		gameDefinitions := requestHandler.FindAll()
		c.JSON(http.StatusOK, gameDefinitions)
	} else {
		gameDefinitions := requestHandler.Find(user.User.ID)
		c.JSON(http.StatusOK, gameDefinitions)
	}

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
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Invalid body content on updateMyGameDefinition")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// dont check ownership when user is a system editor
	skipCheckOwnerShip := auth.UserBelongsToRole(user, auth.SystemEditorRole)

	err = requestHandler.Update(user.User.ID, gameDefinition, skipCheckOwnerShip)
	if err != nil {
		c.AbortWithStatus(http.StatusConflict)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

// updateMyGameDefinition godoc
// @Summary update gamedefition availability by classroom
// @Accept json
// @Produce json
// @Param request body model.GameDefinitionClassroomAvailability true "availability"
// @Success 200 {string} string
// @Security ApiKeyAuth
// @Router /private/mapeditor/update-classroom-map-availability [post]
func updateClassroomMapAvailability(c *gin.Context) {
	user := httphelper.UserDetailsFromContext(c)

	// parse body parameter
	var availability *model.GameDefinitionClassroomAvailability
	err := c.BindJSON(&availability)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Invalid body content on updateClassroomMapAvailability")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// dont check ownership when user is a system editor
	skipCheckOwnerShip := auth.UserBelongsToRole(user, auth.SystemEditorRole)

	err = requestHandler.UpdateAvailability(user.User.ID, availability, skipCheckOwnerShip)
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

// FindAll godoc
func (handler *RequestHandler) FindAll() *[]model.GameDefinition {
	return handler.ds.FindAllSystemGameDefinition()
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
func (handler *RequestHandler) Update(userID uint, gameDefinition *model.GameDefinition, skipCheckOwnerShip bool) error {
	foundByID := handler.ds.FindGameDefinition(gameDefinition.ID)
	// must exist to be updated
	if foundByID == nil {
		log.WithFields(log.Fields{}).Info("gamedefinition DOES NOT EXIST, cant be updated")
		return errors.New("Gamedefinition DOES NOT exist")
	}

	// must be the owner to update it
	if !skipCheckOwnerShip {
		if foundByID.OwnerUserID != userID {
			log.WithFields(log.Fields{
				"foundByID.OwnerUserID": foundByID.OwnerUserID,
				"userID":                userID,
			}).Info("current user dont OWNS this gamedefinition, cant be updated")
			return errors.New("current user DOES NOT OWN this Gamedefinition")
		}
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

// err = requestHandler.UpdateAvailability(user.User.ID, availability, skipCheckOwnerShip)

// Update godoc
func (handler *RequestHandler) UpdateAvailability(userID uint, availability *model.GameDefinitionClassroomAvailability, skipCheckOwnerShip bool) error {
	foundByID := handler.ds.FindGameDefinition(availability.GameDefinitionID)
	// must exist to be updated
	if foundByID == nil {
		log.WithFields(log.Fields{}).Info("gamedefinition DOES NOT EXIST, cant be updated")
		return errors.New("Gamedefinition DOES NOT exist")
	}

	// must be the owner to update it
	if !skipCheckOwnerShip {
		if foundByID.OwnerUserID != userID {
			log.WithFields(log.Fields{
				"foundByID.OwnerUserID": foundByID.OwnerUserID,
				"userID":                userID,
			}).Info("current user dont OWNS this gamedefinition, cant be updated")
			return errors.New("current user DOES NOT OWN this Gamedefinition")
		}
	}

	// load all Available matches for this gamedefinition
	var availableMatches []model.AvailableMatch
	filter := model.AvailableMatch{GameDefinitionID: availability.GameDefinitionID}
	handler.ds.DB.Where(&filter).First(&availableMatches)

	// check records to delete
	for _, search := range availableMatches {
		found := false
		for _, classroom := range availability.Classrooms {
			if search.ClassroomID == classroom {
				found = true
				break
			}
		}

		// this available match was not found in the list
		if !found {
			log.WithFields(log.Fields{
				"availableMatch": search,
			}).Info("remove available match")
			handler.ds.DB.Delete(&search)
		}
	}

	// check records to insert
	for _, classroom := range availability.Classrooms {
		found := false
		for _, availableMatch := range availableMatches {
			if availableMatch.ClassroomID == classroom {
				found = true
				break
			}
		}

		// no available match was found for this classroom, CREATE ONE
		if !found {
			log.WithFields(log.Fields{
				"classroom": classroom,
			}).Info("add missing availability")

			availableMatch := model.AvailableMatch{
				Name:             foundByID.Name,
				ClassroomID:      classroom,
				GameDefinitionID: availability.GameDefinitionID,
			}

			handler.ds.DB.Model(&availableMatch).Create(&availableMatch)
		}
	}

	return nil

}
