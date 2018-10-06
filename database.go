package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/matryer/try.v1"
)

type DBconfig struct {
	dialect  string
	args     string
	host     string
	database string
	user     string
}

type DataSource struct {
	config *DBconfig
	db     *gorm.DB
}

// BuildMysqlConfig creates a DBconfig for Mysql based on environment variables
func BuildMysqlConfig() *DBconfig {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST")
	connection := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v?charset=utf8&parseTime=True&loc=Local", user, password, host, database)

	return &DBconfig{
		dialect:  "mysql",
		host:     host,
		database: database,
		user:     user,
		args:     connection}
}

// NewDataSource creates a DataSource instance
func NewDataSource(config *DBconfig) *DataSource {
	waitTime := 20 * time.Second
	var db *gorm.DB

	err := try.Do(func(attempt int) (bool, error) {
		var err error
		db, err = gorm.Open(config.dialect, config.args)
		log.WithFields(log.Fields{
			"error":    err,
			"host":     config.host,
			"database": config.database,
			"user":     config.user,
		}).Debug("Failed to connect database, will retry")

		if err != nil {
			time.Sleep(waitTime)
		}
		return attempt < 5, err
	})
	if err != nil {
		log.WithFields(log.Fields{
			"host":     config.host,
			"database": config.database,
			"user":     config.user,
		}).Error("Failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&UserSetting{})

	return &DataSource{db: db, config: config}
}

func (ds *DataSource) findUserByEmail(email string) *User {
	var user User

	if ds.db.Where("email = ?", email).First(&user).RecordNotFound() {
		return nil
	}
	return &user
}

func (ds *DataSource) findUserBySession(UUID string) *User {
	var session Session
	var user User

	if ds.db.Where("UUID = ?", UUID).First(&session).RecordNotFound() {
		return nil
	}

	if ds.db.Where("ID = ?", session.UserID).First(&user).RecordNotFound() {
		return nil
	}

	return &user
}

// Create if doesnt exist
func (ds *DataSource) findUserSettingByUser(user *User) *UserSetting {
	var settings UserSetting
	ds.db.Where(&UserSetting{UserID: user.ID}).FirstOrCreate(&settings)
	return &settings
}

func (ds *DataSource) updateUserSetting(settings *UserSetting) *UserSetting {
	var current UserSetting
	if ds.db.First(&current, settings.ID).RecordNotFound() {
		return nil
	}

	current.LastOption = settings.LastOption
	ds.db.Save(&current)

	log.WithFields(log.Fields{
		"settings": current,
	}).Error("User Setting updated")

	return &current
}

func (ds *DataSource) createSession(user *User) *Session {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Errorf("Error creating session UUID: %v", err)
		return nil
	}

	session := Session{UserID: user.ID, UUID: uuid.String()}
	ds.db.Create(&session)

	log.WithFields(log.Fields{
		"user": session.UserID,
		"uuid": session.UUID,
	}).Error("Session created")

	return &session
}
