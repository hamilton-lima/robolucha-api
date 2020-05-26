package datasource

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/model"
)

// FindAllActivities definition
func (ds *DataSource) FindAllActivities() *[]model.Activity {
	var result []model.Activity

	ds.DB.
		Preload("Skill").
		Find(&result)

	log.WithFields(log.Fields{
		"activities": result,
	}).Debug("FindAllActivities")

	return &result
}

// FindLearningObjectiveByName definition
func (ds *DataSource) FindLearningObjectiveByName(name string) *model.LearningObjective {
	var result model.LearningObjective

	if ds.DB.Where(&model.LearningObjective{Name: name}).Find(&result).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"learning objective": result,
	}).Debug("FindLearningObjectiveByName")

	return &result
}

// AddLearningObjective definition
func (ds *DataSource) AddLearningObjective(c *model.LearningObjective) *model.LearningObjective {

	objective := model.LearningObjective{
		Name:   c.Name,
		Skills: c.Skills,
	}

	log.WithFields(log.Fields{
		"objective": objective,
	}).Debug("AddLearningObjective")

	ds.DB.Create(&objective)

	log.WithFields(log.Fields{
		"objective": objective,
	}).Debug("after AddLearningObjective")

	return &objective
}
