package datasource

import (
	"os"
	"testing"

	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/test"
	"gotest.tools/assert"

	log "github.com/sirupsen/logrus"
)

var ds *DataSource

func Setup(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	os.Setenv("GIN_MODE", "release")

	err := os.Remove(test.DB_NAME)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error removing TEST database")
	}

	ds = NewDataSource(BuildSQLLiteConfig(test.DB_NAME))
}

func TestAddGrade(t *testing.T) {
	Setup(t)

	grade := model.Grade{
		Name:    "foo",
		Lowest:  0,
		Highest: 10,
		Color:   "#FFFFFF",
	}

	ds.AddGrade(&grade)
	found := ds.FindGradeByName("foo")
	assert.Equal(t, found.Name, "foo")
}

func TestFindGrade(t *testing.T) {
	Setup(t)
	found := ds.FindGradeByName("orange")
	assert.Assert(t, found == nil)
}
