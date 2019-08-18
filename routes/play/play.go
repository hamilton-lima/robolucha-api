package play

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
)

var requestHandler *RequestHandler
var ds *datasource.DataSource
var publisher pubsub.Publisher

// Init receive database and message queue objects
func Init(_ds *datasource.DataSource, _publisher pubsub.Publisher) *Router {
	ds = _ds
	publisher = _publisher
	requestHandler = Listen()

	return &Router{}
}

// Router definition
type Router struct{}

// Setup definition
func (router *Router) Setup(group *gin.RouterGroup) {
	group.POST("/play", play)
}

// Match definition
type Match struct {
	MatchID uint `json:"matchID"`
}

// play godoc
// @Summary request to play a match
// @Accept json
// @Produce json
// @Param request body model.AvailableMatch true "AvailableMatch"
// @Success 200 {object} play.Match
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

	wait := requestHandler.Send(Request{Data: input})
	response := <-wait
	result := response.Data

	log.WithFields(log.Fields{
		"Match": result,
	}).Info("play()")

	c.JSON(http.StatusOK, result)
}
