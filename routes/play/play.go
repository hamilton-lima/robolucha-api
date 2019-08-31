package play

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/httphelper"
	"gitlab.com/robolucha/robolucha-api/pubsub"
)

var requestHandler *RequestHandler

// Init receive database and message queue objects
func Init(_ds *datasource.DataSource, _publisher pubsub.Publisher) *Router {
	requestHandler = NewRequestHandler(_ds, _publisher)

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
	group.POST("/play/:id", play)
}

// play godoc
// @Summary request to play a match
// @Accept json
// @Produce json
// @Param id path int true "AvailableMatch id"
// @Success 200 {object} model.Match
// @Security ApiKeyAuth
// @Router /private/play/{id} [post]
func play(c *gin.Context) {

	id, err := httphelper.GetIntegerParam(c, "id", "play")
	if err != nil {
		log.Info("Invalid body content on play")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	input := requestHandler.FindAvailableMatchByID(id)
	if input == nil {
		log.WithFields(log.Fields{
			"message": "AvailableMatch not found",
			"id":      id,
		}).Error("Invalid body content on play()")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"AvailableMatch": input,
	}).Info("play()")

	user := httphelper.UserFromContext(c)
	luchador := requestHandler.ds.FindLuchador(user)
	log.WithFields(log.Fields{
		"luchador": luchador,
		"user.id":  user.ID,
	}).Info("publishJoinMatch()")

	match := requestHandler.Play(input, luchador.ID)

	log.WithFields(log.Fields{
		"Match": match,
	}).Info("play()")

	c.JSON(http.StatusOK, match)
}
