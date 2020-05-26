package setup

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"gitlab.com/robolucha/robolucha-api/datasource"
	_ "gitlab.com/robolucha/robolucha-api/docs"
	"gitlab.com/robolucha/robolucha-api/model"
)

// SetupGradeFromFolder definition
func SetupGradeFromFolder(folderName string, ds *datasource.DataSource) {

	files := readFilesFromFolder(folderName)

	for _, file := range files {
		fullPath := filepath.Join(folderName, file.Name())
		log.WithFields(log.Fields{
			"filename": fullPath,
		}).Info("Loading grade")

		CreateGrade(fullPath, ds)
	}

}

// CreateGrade definition
func CreateGrade(fileName string, ds *datasource.DataSource) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.WithFields(log.Fields{
			"fileName": fileName,
			"error":    err,
		}).Error("Error reading grade file")
		return
	}

	jsonContent := string(bytes)
	log.WithFields(log.Fields{
		"jsonContent": jsonContent,
		"filename":    fileName,
	}).Debug("Loading grade")

	var grade model.Grade
	json.Unmarshal(bytes, &grade)

	foundByName := ds.FindGradeByName(grade.Name)
	if foundByName != nil {
		log.WithFields(log.Fields{
			"grade":    grade,
			"filename": fileName,
		}).Info("grade already EXISTS")
	} else {
		ds.AddGrade(&grade)
		log.WithFields(log.Fields{
			"grade":    grade,
			"filename": fileName,
		}).Info("grade CREATED")
	}

}
