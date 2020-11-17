package media

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"bytes"
	"image"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/httphelper"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/gofrs/uuid"
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

func after(text string, find string) string {
	found := strings.LastIndex(text, find)
	if found == -1 {
		return ""
	}

	pos := found + len(find)
	if pos >= len(text) {
		return ""
	}
	return text[pos:len(text)]
}

// Add godoc
func (handler *RequestHandler) AddMedia(request *model.MediaRequest, userID uint) model.Media {

	u2, err := uuid.NewV4()
	if err != nil {
		log.WithFields(log.Fields{
			"step": "error generating UUID",
			"err":  err,
		}).Error("addMedia")
	}

	name := fmt.Sprintf("/tmp/%v-%v", u2, request.FileName)
	thumbnail := fmt.Sprintf("/tmp/%v-thumb-%v", u2, request.FileName)

	// removes "data:image/png;base64," from the beginning of the data
	base64 := after(request.Base64Data, ",")
	data, _ := b64.StdEncoding.DecodeString(base64)

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.WithFields(log.Fields{
			"step": "error creating image",
			"err":  err,
		}).Error("addMedia")
	}

	dstImage800 := imaging.Resize(img, 300, 0, imaging.NearestNeighbor)

	err = ioutil.WriteFile(name, data, 0666)
	if err != nil {
		log.WithFields(log.Fields{
			"step": "error writing temp image file",
			"err":  err,
		}).Error("addMedia")
	}

	err = imaging.Save(dstImage800, thumbnail, imaging.JPEGQuality(85))
	if err != nil {
		log.WithFields(log.Fields{
			"step": "error saving thumbnail temp image file",
			"err":  err,
		}).Error("addMedia")
	}

	// upload files
	uploadOriginal, errOriginal := upload(name)
	if errOriginal != nil {
		log.WithFields(log.Fields{
			"step": "error uploading original file",
			"err":  errOriginal,
		}).Error("addMedia")
	}

	uploadThumbnail, errThumb := upload(thumbnail)
	if errThumb != nil {
		log.WithFields(log.Fields{
			"step": "error uploading thumbnail",
			"err":  errThumb,
		}).Error("addMedia")
	}

	err = os.Remove(name)
	if err != nil {
		log.WithFields(log.Fields{
			"step": "error removing temp file",
			"name": name,
			"err":  err,
		}).Error("addMedia")
	}

	err = os.Remove(thumbnail)
	if err != nil {
		log.WithFields(log.Fields{
			"step": "error removing thumbnail temp file",
			"name": thumbnail,
			"err":  err,
		}).Error("addMedia")
	}

	log.WithFields(log.Fields{
		"step":            "after upload",
		"uploadThumbnail": uploadThumbnail,
		"uploadOriginal":  uploadOriginal,
	}).Info("addMedia")

	// upload the file here
	media := model.Media{
		UserID:    userID,
		FileName:  request.FileName,
		URL:       uploadOriginal.Location,
		Thumbnail: uploadThumbnail.Location,
	}

	handler.ds.DB.Create(&media)
	return media
}

func upload(fileName string) (*s3manager.UploadOutput, error) {
	// https://jto.nyc3.digitaloceanspaces.com
	// The session the S3 Uploader will use
	endpoint := "nyc3.digitaloceanspaces.com"
	region := "nyc3"
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
	}))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q, %v", fileName, err)
	}

	myBucket := "game-robolucha"
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(fileName),
		Body:   f,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file, %v", err)
	}

	return result, err
}
