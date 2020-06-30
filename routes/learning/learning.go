package learning

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/httphelper"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/robolucha/robolucha-api/datasource"
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
	group.GET("/activity", getActivity)
	group.GET("/assignment", getAssignments)
	group.GET("/assignment/:id", getAssignment)
	group.POST("/assignment", addAssignment)
	group.DELETE("/assignment/:id", delAssignment)
	group.PATCH("/assignment/:id/students", updateAssignmentStudents)
	group.PATCH("/assignment/:id/activities", updateAssignmentActivities)
	group.GET("/badword", checkBadWord)

}

func checkBadWord(context *gin.Context) {
	sentence := context.Query("sentence")
	context.JSON(http.StatusOK, utility.ContainsBadWord(sentence))
}

// updateAssignmentActivities godoc
// @Summary find existing activities
// @Accept json
// @Produce json
// @Success 200 {array} int
// @Security ApiKeyAuth
// @Router /assignment/:id/activities [patch]
func updateAssignmentActivities(c *gin.Context) {
	var activityIds []uint
	id, _ := httphelper.GetIntegerParam(c, "id", "updateAssignmentActivities")
	c.BindJSON(&activityIds)
	result := requestHandler.ds.UpdateAssignmentActivities(id, activityIds)
	c.JSON(http.StatusOK, result)
}

// updateAssignmentStudents godoc
// @Summary find existing activities
// @Accept json
// @Produce json
// @Success 200 {array} int
// @Security ApiKeyAuth
// @Router /assignment/:id/students [patch]
func updateAssignmentStudents(c *gin.Context) {
	var studentIds []uint
	id, _ := httphelper.GetIntegerParam(c, "id", "updateAssignmentStudents")
	c.BindJSON(&studentIds)
	result := requestHandler.ds.UpdateAssignmentStudents(id, studentIds)
	c.JSON(http.StatusOK, result)
}


// getActivity godoc
// @Summary find existing activities
// @Accept json
// @Produce json
// @Success 200 {array} model.Activity
// @Security ApiKeyAuth
// @Router /dashboard/activity [get]
func getActivity(c *gin.Context) {
	result := requestHandler.ds.FindAllActivities()

	log.WithFields(log.Fields{
		"activities": result,
	}).Info("getActivity")

	c.JSON(http.StatusOK, result)
}

// getAssignments godoc
// @Summary find existing activities
// @Accept json
// @Produce json
// @Success 200 {array} model.Activity
// @Security ApiKeyAuth
// @Router /dashboard/assignments [get]
func getAssignments(c *gin.Context) {
	result := requestHandler.ds.FindAllAssignments()

	//log.WithFields(log.Fields{
	//	"activities": result,
	//}).Info("getActivity")

	c.JSON(http.StatusOK, result)
}

// getAssignment godoc
// @Summary find existing activities
// @Accept json
// @Produce json
// @Success 200 {array} model.Activity
// @Security ApiKeyAuth
// @Router /dashboard/assignments [get]
func getAssignment(c *gin.Context) {
	id, _ := httphelper.GetIntegerParam(c, "id", "getAssignment")
	result := requestHandler.ds.FindAssignmentById(id)

	c.JSON(http.StatusOK, result)
}

// addAssignment godoc
// @Summary find existing activities
// @Accept json
// @Produce json
// @Success 200 {array} model.Activity
// @Security ApiKeyAuth
// @Router /dashboard/assignments [post]
func addAssignment(c *gin.Context) {
	var assignment *model.Assignment
	c.BindJSON(&assignment)
	result := requestHandler.ds.AddAssignment(assignment)
	c.JSON(http.StatusOK, result)
}

// delAssignment godoc
// @Summary find existing activities
// @Accept json
// @Produce json
// @Success 200 {array} model.Activity
// @Security ApiKeyAuth
// @Router /dashboard/assignments [get]
func delAssignment(c *gin.Context) {
	id, _ := httphelper.GetIntegerParam(c, "id", "getAssignment")
	requestHandler.ds.DeleteAssignment(id)
	c.JSON(http.StatusOK, id)
}

