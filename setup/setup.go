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

// LoadMetadataFromFolder loads all metadata from folderName
func LoadMetadataFromFolder(folderName string, ds *datasource.DataSource) {
	SetupGameDefinitionFromFolder(filepath.Join(folderName, "gamedefinition"), ds)
	SetupGradeFromFolder(filepath.Join(folderName, "grade"), ds)
}

func SetupGameDefinitionFromFolder(folderName string, ds *datasource.DataSource) {

	log.WithFields(log.Fields{
		"folderName": folderName,
	}).Info("SetupGameDefinitionFromFolder")

	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		log.WithFields(log.Fields{
			"folderName": folderName,
			"error":      err,
		}).Error("Error loading gamedefinition files")
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(folderName, file.Name())
		log.WithFields(log.Fields{
			"filename": fullPath,
		}).Info("Loading gamedefinition")

		CreateGameDefinition(fullPath, ds)
	}

}

// CreateGameDefinition definition
func CreateGameDefinition(fileName string, ds *datasource.DataSource) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.WithFields(log.Fields{
			"fileName": fileName,
			"error":    err,
		}).Error("Error reading gamedefinition file")
		return
	}

	jsonContent := string(bytes)
	log.WithFields(log.Fields{
		"jsonContent": jsonContent,
		"filename":    fileName,
	}).Debug("Loading gamedefinition")

	var gameDefinition model.GameDefinition
	json.Unmarshal(bytes, &gameDefinition)

	foundByName := ds.FindGameDefinitionByName(gameDefinition.Name)
	if foundByName != nil {
		log.WithFields(log.Fields{
			"gameDefinition": gameDefinition,
			"filename":       fileName,
		}).Info("gamedefinition already EXISTS")
	} else {
		createResult := ds.CreateGameDefinition(&gameDefinition)
		log.WithFields(log.Fields{
			"gameDefinition": createResult,
			"filename":       fileName,
		}).Info("gamedefinition CREATED")
	}

}

func CreateAvailableMatches(ds *datasource.DataSource) {

	all := ds.FindAllGameDefinition()
	for _, gd := range *all {
		ds.CreateAvailableMatchIfDontExist(gd.ID, gd.Name)
	}

}

// SetupGradeFromFolder definition
func SetupGradeFromFolder(folderName string, ds *datasource.DataSource) {

	log.WithFields(log.Fields{
		"folderName": folderName,
	}).Info("SetupGradeFromFolder")

	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		log.WithFields(log.Fields{
			"folderName": folderName,
			"error":      err,
		}).Error("Error loading grade files")
		return
	}

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
