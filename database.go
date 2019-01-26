package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
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
		}).Debug("Database connection status")

		if err != nil {
			log.WithFields(log.Fields{
				"waitTime": waitTime,
			}).Warn("Error connecting to the database, will retry.")

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
	db.AutoMigrate(&Match{})
	db.AutoMigrate(&Luchador{})
	db.AutoMigrate(&Code{})
	db.AutoMigrate(&Config{})
	db.AutoMigrate(&MatchParticipant{})

	// Enable debug mode
	debug := os.Getenv("GORM_DEBUG")
	if debug == "true" {
		db.LogMode(true)
	}

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

func (ds *DataSource) createMatch(m *Match) *Match {
	match := Match{
		TimeStart:     m.TimeStart,
		TimeEnd:       m.TimeEnd,
		LastTimeAlive: m.LastTimeAlive,
		Duration:      m.Duration,
	}

	ds.db.Create(&match)

	log.WithFields(log.Fields{
		"match.id": match.ID,
		"duration": match.Duration,
	}).Info("Match created")

	return &match
}

func (ds *DataSource) createLuchador(l *Luchador) *Luchador {
	luchador := Luchador{
		UserID:  l.UserID,
		Name:    l.Name,
		Codes:   l.Codes,
		Configs: l.Configs,
	}

	ds.db.Create(&luchador)

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("Luchador created")

	return &luchador
}

func (ds *DataSource) findLuchador(user *User) *Luchador {
	var luchador Luchador
	if ds.db.Preload("Codes").Preload("Configs").Where(&Luchador{UserID: user.ID}).First(&luchador).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchador")

	return &luchador
}

func (ds *DataSource) updateLuchador(user *User, luchador *Luchador) *Luchador {
	var current Luchador
	if ds.db.First(&current, luchador.ID).RecordNotFound() {
		return nil
	}

	ds.db.Save(&luchador)

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("updateLuchador")

	return luchador
}

func (ds *DataSource) findActiveMatches() *[]Match {

	var matches []Match
	ds.db.Where("time_end < time_start").Order("time_start desc").Find(&matches)

	log.WithFields(log.Fields{
		"matches": matches,
	}).Info("findActiveMatches")

	return &matches
}

func (ds *DataSource) findMatch(id uint) *Match {

	var match Match
	ds.db.Where(&Match{ID: id}).Find(&match)

	log.WithFields(log.Fields{
		"id":    id,
		"match": match,
	}).Info("findMatch")

	return &match
}

func (ds *DataSource) findLuchadorByID(luchadorID uint) *Luchador {
	var luchador Luchador
	if ds.db.Preload("Codes").Preload("Configs").Where(&Luchador{ID: luchadorID}).First(&luchador).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchadorByID")

	return &luchador
}

func (ds *DataSource) createMatchParticipant(mp *MatchParticipant) *MatchParticipant {
	matchPartipant := MatchParticipant{
		LuchadorID: mp.LuchadorID,
		MatchID:    mp.MatchID,
	}

	ds.db.Create(&matchPartipant)

	log.WithFields(log.Fields{
		"matchPartipant": matchPartipant,
	}).Info("MatchPartipant created")

	return &matchPartipant
}

func (ds *DataSource) endMatch(match *Match) *Match {

	ds.db.Model(&match).Update("time_end", match.TimeEnd)

	log.WithFields(log.Fields{
		"match": match,
	}).Info("Match time_end updated")

	return match
}
