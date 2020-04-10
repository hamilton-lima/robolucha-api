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
