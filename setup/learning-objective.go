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

// SetupLearningObjectiveFromFolder definition
func SetupLearningObjectiveFromFolder(folderName string, ds *datasource.DataSource) {

	files := readFilesFromFolder(folderName)

	for _, file := range files {
		fullPath := filepath.Join(folderName, file.Name())
		log.WithFields(log.Fields{
			"filename": fullPath,
		}).Info("Loading learning objective")

		CreateLearningObjective(fullPath, ds)
	}

}

// CreateLearningObjective definition
func CreateLearningObjective(fileName string, ds *datasource.DataSource) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.WithFields(log.Fields{
			"fileName": fileName,
			"error":    err,
		}).Error("Error reading learning objective file")
		return
	}

	jsonContent := string(bytes)
	log.WithFields(log.Fields{
		"jsonContent": jsonContent,
		"filename":    fileName,
	}).Debug("Loading learning objective")

	var objective model.LearningObjective
	json.Unmarshal(bytes, &objective)

	foundByName := ds.FindLearningObjectiveByName(objective.Name)
	if foundByName != nil {
		log.WithFields(log.Fields{
			"learning objective": objective,
			"filename":           fileName,
		}).Info("learning objective already EXISTS")
	} else {
		ds.AddLearningObjective(&objective)
		log.WithFields(log.Fields{
			"learning objective": objective,
			"filename":           fileName,
		}).Info("learning objective CREATED")
	}

}
