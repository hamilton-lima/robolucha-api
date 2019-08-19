package play

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/httphelper"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
)

var requestHandler *RequestHandler

// Init receive database and message queue objects
func Init(_ds *datasource.DataSource, _publisher pubsub.Publisher) *Router {
	requestHandler = Listen(_ds, _publisher)

	return &Router{ds: _ds,
		publisher: _publisher,
	}
}

// Router definition
type Router struct {
	ds        *datasource.DataSource
	publisher pubsub.Publisher
}

// Setup definition
func (router *Router) Setup(group *gin.RouterGroup) {
	group.POST("/play", play)
}

// play godoc
// @Summary request to play a match
// @Accept json
// @Produce json
// @Param request body model.AvailableMatch true "AvailableMatch"
// @Success 200 {object} model.Match
// @Security ApiKeyAuth
// @Router /internal/play [post]
func play(c *gin.Context) {

	var input *model.AvailableMatch
	err := c.BindJSON(&input)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Invalid body content on play()")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"AvailableMatch": input,
	}).Info("play()")

	user := httphelper.UserFromContext(c)

	wait := requestHandler.Send(Request{User: user, AvailableMatch: input})
	response := <-wait

	log.WithFields(log.Fields{
		"Match": response.Match,
	}).Info("play()")

	c.JSON(http.StatusOK, response.Match)
}
