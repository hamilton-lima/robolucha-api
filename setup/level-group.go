package setup

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
)

// SetupLevelGroupFromFolder definition
func SetupLevelGroupFromFolder(folderName string, ds *datasource.DataSource) {

	files := readFilesFromFolder(folderName)

	for _, file := range files {
		fullPath := filepath.Join(folderName, file.Name())
		log.WithFields(log.Fields{
			"filename": fullPath,
		}).Info("Loading level group")

		CreateLevelGroup(fullPath, ds)
	}

}

// CreateLevelGroup definition
func CreateLevelGroup(fileName string, ds *datasource.DataSource) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.WithFields(log.Fields{
			"fileName": fileName,
			"error":    err,
		}).Error("Error reading level group file")
		return
	}

	jsonContent := string(bytes)
	log.WithFields(log.Fields{
		"jsonContent": jsonContent,
		"filename":    fileName,
	}).Debug("Loading level group")

	var levelGroup model.LevelGroup
	json.Unmarshal(bytes, &levelGroup)

	foundByName := ds.FindLevelGroupByName(levelGroup.Name)
	if foundByName != nil {
		log.WithFields(log.Fields{
			"levelGroup": levelGroup,
			"filename":   fileName,
		}).Info("level group already EXISTS")
	} else {
		ds.AddLevelGroup(&levelGroup)
		log.WithFields(log.Fields{
			"levelGroup": levelGroup,
			"filename":   fileName,
		}).Info("level group CREATED")
	}
}
