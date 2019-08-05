package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	_ "gitlab.com/robolucha/robolucha-api/docs"
	"gitlab.com/robolucha/robolucha-api/model"
)

func SetupGameDefinitionFromFolder(folderName string) {

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

		CreateGameDefinition(fullPath)
	}

}

// CreateGameDefinition definition
func CreateGameDefinition(fileName string) {
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

	foundByName := dataSource.findGameDefinitionByName(gameDefinition.Name)
	if foundByName != nil {
		log.WithFields(log.Fields{
			"gameDefinition": gameDefinition,
			"filename":       fileName,
		}).Info("gamedefinition already EXISTS")
	} else {
		createResult := dataSource.createGameDefinition(&gameDefinition)
		log.WithFields(log.Fields{
			"gameDefinition": createResult,
			"filename":       fileName,
		}).Info("gamedefinition CREATED")
	}

}
