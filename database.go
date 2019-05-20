package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	secret string
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

// BuildMysqlConfig creates a DBconfig for Mysql based on environment variables
func BuildSQLLiteConfig(fileName string) *DBconfig {
	return &DBconfig{
		dialect: "sqlite3",
		args:    fileName}
}

// NewDataSource creates a DataSource instance
func NewDataSource(config *DBconfig) *DataSource {
	waitTime := 2 * time.Second
	var db *gorm.DB

	log.WithFields(log.Fields{
		"host":     config.host,
		"database": config.database,
		"user":     config.user,
	}).Debug("Connecting to the database")

	err := try.Do(func(attempt int) (bool, error) {
		var err error
		db, err = gorm.Open(config.dialect, config.args)

		// Enable debug mode
		debug := os.Getenv("GORM_DEBUG")
		if debug == "true" {
			db.LogMode(true)
		}

		log.WithFields(log.Fields{
			"error":    err,
			"host":     config.host,
			"database": config.database,
			"user":     config.user,
		}).Debug("Database connection status")

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
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&UserSetting{})
	db.AutoMigrate(&Match{})
	db.AutoMigrate(&Luchador{})
	db.AutoMigrate(&Code{})
	db.AutoMigrate(&Config{})
	db.AutoMigrate(&MatchScore{})
	db.AutoMigrate(&ServerCode{})
	db.AutoMigrate(&SceneComponent{})
	db.AutoMigrate(&GameDefinition{})

	secret := os.Getenv("API_SECRET")

	return &DataSource{db: db, config: config, secret: secret}
}

func (ds *DataSource) KeepAlive() {
	log.Debug("Keep connection alive")
	for range time.Tick(time.Minute) {
		ds.db.DB().Ping()
		log.Debug("Keep connection alive")
	}
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

func (ds *DataSource) createUser(u User) *User {
	user := User{Username: u.Username}
	ds.db.Where(&User{Username: u.Username}).FirstOrCreate(&user)

	log.WithFields(log.Fields{
		"id":       user.ID,
		"username": user.Username,
	}).Debug("createUser")

	return &user
}

func (ds *DataSource) createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key + ds.secret))
	return hex.EncodeToString(hasher.Sum(nil))
}

// func (ds *DataSource) createSession(user *User) *Session {
// 	uuid, err := uuid.NewV4()
// 	if err != nil {
// 		log.Errorf("Error creating session UUID: %v", err)
// 		return nil
// 	}

// 	session := Session{UserID: user.ID, UUID: uuid.String()}
// 	ds.db.Create(&session)

// 	log.WithFields(log.Fields{
// 		"user": session.UserID,
// 		"uuid": session.UUID,
// 	}).Info("Session created")

// 	return &session
// }

func (ds *DataSource) createMatch(m *Match) *Match {
	match := Match{
		TimeStart:     m.TimeStart,
		TimeEnd:       m.TimeEnd,
		LastTimeAlive: m.LastTimeAlive,
		Duration:      m.Duration,
		Participants:  m.Participants,
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

func (ds *DataSource) findLuchadorByIDNoPreload(id uint) *Luchador {
	var luchador Luchador
	if ds.db.First(&luchador, id).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchadorByID")

	return &luchador
}

func (ds *DataSource) updateLuchador(luchador *Luchador) *Luchador {
	current := ds.findLuchadorByID(luchador.ID)
	if current == nil {
		return nil
	}

	current.Name = luchador.Name
	current.Configs = applyConfigChanges(current.Configs, luchador.Configs)
	current.Codes = luchador.Codes

	ds.db.Save(current)

	log.WithFields(log.Fields{
		"luchador": current,
	}).Info("after updateLuchador")

	return current
}

func applyConfigChanges(original []Config, updated []Config) []Config {
	for i, configOriginal := range original {
		for _, configUpdated := range updated {
			if configOriginal.Key == configUpdated.Key {
				// NOTE that range make COPIES of the values!!
				original[i].Value = configUpdated.Value
				break
			}
		}
	}
	return original
}

func (ds *DataSource) findActiveMatches() *[]Match {

	var matches []Match
	ds.db.Where("time_end < time_start").Order("time_start desc").Find(&matches)

	log.WithFields(log.Fields{
		"matches": matches,
	}).Info("findActiveMatches")

	return &matches
}

func (ds *DataSource) findMaskConfig(id uint) *[]Config {

	var configs []Config
	ds.db.Where(&Config{LuchadorID: id}).Find(&configs)

	log.WithFields(log.Fields{
		"configs": configs,
	}).Info("findMaskConfig")

	return &configs
}

func (ds *DataSource) findMatch(id uint) *Match {

	var match Match
	ds.db.Preload("Participants").Where(&Match{ID: id}).First(&match)

	log.WithFields(log.Fields{
		"id":    id,
		"match": match,
	}).Info("findMatch")

	return &match
}

// func (ds *DataSource) findGameComponentByID(id uint) *Luchador {
// 	var luchador Luchador
// 	if ds.db.First(&luchador, id).RecordNotFound(){
// 		return nil
// 	}

// 	log.WithFields(log.Fields{
// 		"luchador": luchador,
// 	}).Info("findGameComponentByID")

// 	return &gameComponent
// }

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

func (ds *DataSource) findLuchadorByName(name string) *Luchador {
	var luchador Luchador
	if ds.db.Where(&Luchador{Name: name}).First(&luchador).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchadorByName")

	return &luchador
}

func (ds *DataSource) NameExist(ID uint, name string) bool {
	var luchador Luchador
	result := !ds.db.Where("id <> ? AND name = ?", ID, name).First(&luchador).RecordNotFound()

	log.WithFields(log.Fields{
		"luchador": luchador,
		"result":   result,
	}).Debug("NameExist")

	return result
}

func (ds *DataSource) addMatchParticipant(mp *MatchParticipant) *MatchParticipant {

	var match *Match
	match = ds.findMatch(mp.MatchID)
	if match == nil {
		log.WithFields(log.Fields{
			"matchID": mp.MatchID,
		}).Error("Match not found")
		return nil
	}

	var luchador *Luchador
	luchador = ds.findLuchadorByIDNoPreload(mp.LuchadorID)
	if luchador == nil {
		log.WithFields(log.Fields{
			"luchadorID": mp.LuchadorID,
		}).Error("Luchador not found")
		return nil
	}

	for _, participant := range match.Participants {
		if participant.ID == mp.LuchadorID {
			log.WithFields(log.Fields{
				"matchID":    mp.MatchID,
				"luchadorID": mp.LuchadorID,
			}).Error("Luchador is already in the match")
			return nil
		}
	}

	match.Participants = append(match.Participants, *luchador)
	ds.db.Save(&match)

	matchPartipant := MatchParticipant{
		LuchadorID: luchador.ID,
		MatchID:    match.ID,
	}

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

func (ds *DataSource) findLuchadorConfigsByMatchID(id uint) *[]Luchador {

	match := Match{}
	ds.db.First(&match, "id = ?", id)

	var participants []Luchador
	ds.db.Model(&match).Related(&participants, "Participants").Preload("Configs")

	log.WithFields(log.Fields{
		"id":     id,
		"match":  match,
		"result": participants,
	}).Debug("findLuchadorConfigsByMatchID")

	return &participants
}

func (ds *DataSource) getMatchScoresByMatchID(id uint) *[]MatchScore {

	result := []MatchScore{}
	ds.db.Where(&MatchScore{MatchID: id}).Find(&result)

	for _, val := range result {
		log.WithFields(log.Fields{
			"matchId":    val.MatchID,
			"luchadorId": val.LuchadorID,
			"score":      val.Score,
		}).Debug("getMatchScoresByMatchID")
	}

	return &result
}

func (ds *DataSource) addMatchScores(ms *ScoreList) *ScoreList {

	log.WithFields(log.Fields{
		"action":    "start",
		"scorelist": ms,
	}).Info("addMatchScores")

	var match *Match = nil

	for _, score := range ms.Scores {

		if match == nil {
			match = ds.findMatch(score.MatchID)
			if match == nil {
				log.WithFields(log.Fields{
					"matchID": score.MatchID,
				}).Error("Match not found")
				return nil
			}

			log.WithFields(log.Fields{
				"action": "match-found",
				"match":  match,
			}).Info("addMatchScores")
		}

		var luchador *Luchador
		luchador = ds.findLuchadorByID(score.LuchadorID)
		if luchador == nil {
			log.WithFields(log.Fields{
				"luchadorID": score.LuchadorID,
			}).Error("Luchador not found")
			return nil
		}

		log.WithFields(log.Fields{
			"action":   "luchador-found",
			"luchador": luchador,
		}).Info("addMatchScores")

		score := MatchScore{
			LuchadorID: luchador.ID,
			MatchID:    match.ID,
			Kills:      score.Kills,
			Deaths:     score.Deaths,
			Score:      score.Score,
		}

		log.WithFields(log.Fields{
			"action": "before-save",
			"score":  score,
		}).Info("addMatchScores")

		ds.db.Create(&score)

		log.WithFields(log.Fields{
			"action": "after-save",
			"score":  score,
		}).Info("addMatchScores")
	}

	return ms
}

func (ds *DataSource) createGameDefinition(g *GameDefinition) *GameDefinition {

	gameDefinition := GameDefinition{}
	copier.Copy(&gameDefinition, &g)
	ds.db.Create(&gameDefinition)

	log.WithFields(log.Fields{
		"gameDefinition": gameDefinition,
	}).Info("createGameDefinition")

	return &gameDefinition
}

func (ds *DataSource) findGameDefinition(id uint) *GameDefinition {
	var gameDefinition GameDefinition

	if ds.db.
		Preload("Participants").
		Preload("SceneComponents").
		Preload("Codes").
		Preload("LuchadorSuggestedCodes").
		Where(&GameDefinition{ID: id}).
		First(&gameDefinition).
		RecordNotFound() {

		log.WithFields(log.Fields{
			"ID": id,
		}).Info("findGameDefinition not found")

		return nil
	}

	log.WithFields(log.Fields{
		"ID":             id,
		"gameDefinition": gameDefinition,
	}).Info("findGameDefinition before array checks")

	if gameDefinition.Participants == nil {
		gameDefinition.Participants = make([]Luchador, 0)
	}

	if gameDefinition.SceneComponents == nil {
		gameDefinition.SceneComponents = make([]SceneComponent, 0)
	}

	if gameDefinition.Codes == nil {
		gameDefinition.Codes = make([]ServerCode, 0)
	}

	if gameDefinition.LuchadorSuggestedCodes == nil {
		gameDefinition.LuchadorSuggestedCodes = make([]ServerCode, 0)
	}

	log.WithFields(log.Fields{
		"ID":             id,
		"gameDefinition": gameDefinition,
	}).Info("findGameDefinition")

	return &gameDefinition
}
