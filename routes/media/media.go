package media

import (
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
	group.POST("/media", addMedia)
}

// addMedia godoc
// @Summary add media
// @Accept json
// @Produce json
// @Param request body model.MediaRequest true "MediaRequest"
// @Success 200 {object} model.Media
// @Security ApiKeyAuth
// @Router /private/media [post]
func addMedia(c *gin.Context) {
	user := httphelper.UserFromContext(c)

	var request *model.MediaRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Info("Invalid body content on addMedia")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response := requestHandler.AddMedia(request, user.ID)
	log.WithFields(log.Fields{
		"response": response,
	}).Info("addMedia")

	c.JSON(http.StatusOK, response)
}

// Add godoc
func (handler *RequestHandler) AddMedia(request *model.MediaRequest, userID uint) model.Media {

	// upload the file here
	media := model.Media{
		UserID:    userID,
		FileName:  request.FileName,
		URL:       "",
		Thumbnail: "",
	}

	handler.ds.DB.Create(&media)
	return media
}
