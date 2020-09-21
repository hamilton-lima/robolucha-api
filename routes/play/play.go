package play

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/httphelper"
	"gitlab.com/robolucha/robolucha-api/model"
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
	group.POST("/play", play)
	group.POST("/leave-tutorial-match", leaveTutorialMatch)
}

// play godoc
// @Summary request to play a match
// @Accept json
// @Produce json
// @Param request body model.PlayRequest true "PlayRequest"
// @Success 200 {object} model.Match
// @Security ApiKeyAuth
// @Router /private/play [post]
func play(c *gin.Context) {

	var playRequest *model.PlayRequest
	err := c.BindJSON(&playRequest)
	if err != nil {
		log.Info("Invalid body content on play")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	input := requestHandler.FindAvailableMatchByID(playRequest.AvailableMatchID)
	if input == nil {
		log.WithFields(log.Fields{
			"message": "AvailableMatch not found",
			"id":      playRequest.AvailableMatchID,
		}).Error("Invalid body content on play()")

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"AvailableMatch": input,
	}).Info("play()")

	user := httphelper.UserDetailsFromContext(c)

	if requestHandler.UserHasLevelToPlay(&user.Level, input.GameDefinition) {

		luchador := requestHandler.ds.FindLuchador(user.User)
		log.WithFields(log.Fields{
			"luchador": model.LogGameComponent(luchador),
			"user.id":  user.User.ID,
			"TeamID":   playRequest.TeamID,
		}).Info("play()")

		match := requestHandler.Play(input, luchador.ID, playRequest.TeamID)

		log.WithFields(log.Fields{
			"Match": match,
		}).Info("play()")

		c.JSON(http.StatusOK, match)
	} else {
		log.WithFields(log.Fields{
			"user.level":     user.Level.Level,
			"match minlevel": input.GameDefinition.MinLevel,
			"match maxlevel": input.GameDefinition.MaxLevel,
		}).Error("play() user DO NOT have right level to play")

		c.JSON(http.StatusBadRequest, nil)
	}

}

// leaveTutorialMatch godoc
// @Summary Sends message to end active tutorial matches
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Security ApiKeyAuth
// @Router /private/leave-tutorial-match [post]
func leaveTutorialMatch(c *gin.Context) {

	user := httphelper.UserFromContext(c)

	var luchador *model.GameComponent
	luchador = requestHandler.ds.FindLuchador(user)
	if luchador == nil {
		log.WithFields(log.Fields{
			"user": user,
		}).Error("Error getting luchador for the current user")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	requestHandler.LeaveTutorialMatches(luchador)
	c.JSON(http.StatusOK, "")
}
