package setup

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gitlab.com/robolucha/robolucha-api/utility"

	log "github.com/sirupsen/logrus"

	"gitlab.com/robolucha/robolucha-api/datasource"
	_ "gitlab.com/robolucha/robolucha-api/docs"
)

// LoadMetadataFromFolder loads all metadata from folderName
func LoadMetadataFromFolder(folderName string, ds *datasource.DataSource) {
	SetupGameDefinitionFromFolder(filepath.Join(folderName, "gamedefinition"), ds)
	SetupGradeFromFolder(filepath.Join(folderName, "grade"), ds)
	SetupLearningObjectiveFromFolder(filepath.Join(folderName, "learning-objective"), ds)
	SetupLevelGroupFromFolder(filepath.Join(folderName, "level-group"), ds)
	utility.SetupBadWordListFromFolder(folderName)
}

// CreateAvailableMatches definition
func CreateAvailableMatches(ds *datasource.DataSource) {

	all := ds.FindAllSystemGameDefinition()
	for _, gd := range *all {
		ds.CreateAvailableMatchIfDontExist(gd.ID, gd.Name)
	}

}

func readFilesFromFolder(folderName string) []os.FileInfo {

	log.WithFields(log.Fields{
		"folderName": folderName,
	}).Info("readFilesFromFolder")

	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		log.WithFields(log.Fields{
			"folderName": folderName,
			"error":      err,
		}).Error("Error loading files from folder")
		nofiles := make([]os.FileInfo, 0)
		return nofiles
	}

	return files
}
