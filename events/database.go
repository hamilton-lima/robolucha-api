package events

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/model"
	try "gopkg.in/matryer/try.v1"
)

// DBconfig defines database configuration
type DBconfig struct {
	dialect  string
	args     string
	host     string
	database string
	user     string
}

// DataSource keep the connnection instance and the configuration
type DataSource struct {
	config *DBconfig
	DB     *gorm.DB
}

// BuildMysqlConfig creates a DBconfig for Mysql based on environment variables
func BuildMysqlConfig() *DBconfig {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE_EVENTS")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")

	connection := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", user, password, host, port, database)

	return &DBconfig{
		dialect:  "mysql",
		host:     host,
		database: database,
		user:     user,
		args:     connection}
}

// NewDataSource creates a DataSource instance
func NewDataSource(config *DBconfig) *DataSource {
	waitTime := 2 * time.Second
	var DB *gorm.DB

	log.WithFields(log.Fields{
		"host":     config.host,
		"database": config.database,
		"user":     config.user,
	}).Debug("Connecting to the database")

	err := try.Do(func(attempt int) (bool, error) {
		var err error
		DB, err = gorm.Open(config.dialect, config.args)

		// Enable debug mode
		debug := os.Getenv("GORM_DEBUG")
		if debug == "true" {
			DB.LogMode(true)
		}

		log.WithFields(log.Fields{
			"error":    err,
			"host":     config.host,
			"database": config.database,
			"user":     config.user,
		}).Info("Database connection status")

		if err != nil {
			log.WithFields(log.Fields{
				"waitTime": waitTime,
				"err":      err,
			}).Warn("Error connecting to the database, will retry.")

			time.Sleep(waitTime)
		}
		return attempt < 30, err
	})
	if err != nil {
		log.WithFields(log.Fields{
			"host":     config.host,
			"database": config.database,
			"user":     config.user,
		}).Error("Failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(&model.PageEvent{})

	return &DataSource{DB: DB, config: config}
}

// KeepAlive sends ticks to the DB to keep the connection alive
func (ds *DataSource) KeepAlive() {
	log.Debug("Keep connection alive")
	for range time.Tick(time.Minute) {
		ds.DB.DB().Ping()
		log.Debug("Keep connection alive")
	}
}

func (ds *DataSource) CreateEvent(event model.PageEvent) {
	ds.DB.Create(&event)

	log.WithFields(log.Fields{
		"event": event,
	}).Debug("event created")

}
