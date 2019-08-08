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
	db     *gorm.DB
	secret string
}

const GAMEDEFINITION_TYPE_TUTORIAL = "tutorial"
const GAMEDEFINITION_TYPE_MULTIPLAYER = "multiplayer"

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

// BuildSQLLiteConfig creates a DBconfig for Mysql based on environment variables
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
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Session{})
	db.AutoMigrate(&model.UserSetting{})
	db.AutoMigrate(&model.Match{})
	db.AutoMigrate(&model.Code{})
	db.AutoMigrate(&model.Config{})
	db.AutoMigrate(&model.MatchScore{})
	db.AutoMigrate(&model.SceneComponent{})
	db.AutoMigrate(&model.GameComponent{})
	db.AutoMigrate(&model.GameDefinition{})
	db.AutoMigrate(&model.MatchMetric{})
	db.AutoMigrate(&model.Classroom{})
	db.AutoMigrate(&model.Student{})
	db.AutoMigrate(&model.AvailableMatch{})

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

func (ds *DataSource) findUserByEmail(email string) *model.User {
	var user model.User

	if ds.db.Where("email = ?", email).First(&user).RecordNotFound() {
		return nil
	}
	return &user
}

func (ds *DataSource) findUserBySession(UUID string) *model.User {
	var session model.Session
	var user model.User

	if ds.db.Where("UUID = ?", UUID).First(&session).RecordNotFound() {
		return nil
	}

	if ds.db.Where("ID = ?", session.UserID).First(&user).RecordNotFound() {
		return nil
	}

	return &user
}

// Create if doesnt exist
func (ds *DataSource) findUserSettingByUser(user *model.User) *model.UserSetting {
	var settings model.UserSetting
	ds.db.Where(&model.UserSetting{UserID: user.ID}).FirstOrCreate(&settings)
	return &settings
}

func (ds *DataSource) updateUserSetting(settings *model.UserSetting) *model.UserSetting {
	var current model.UserSetting
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

func (ds *DataSource) createUser(name string) *model.User {
	user := model.User{Username: name}
	ds.db.Where(&model.User{Username: name}).FirstOrCreate(&user)

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

func (ds *DataSource) createMatch(gameDefinitionID uint) *model.Match {
	match := model.Match{
		TimeStart: time.Now(),
	}

	gameDefinition := ds.findGameDefinition(gameDefinitionID)
	copier.Copy(&match, &gameDefinition)

	ds.db.Create(&match)

	log.WithFields(log.Fields{
		"match.id":         match.ID,
		"gameDefinitionID": gameDefinitionID,
	}).Info("Match created")

	return &match
}

func (ds *DataSource) createLuchador(l *model.GameComponent) *model.GameComponent {
	luchador := model.GameComponent{
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

func (ds *DataSource) findLuchador(user *model.User) *model.GameComponent {
	var luchador model.GameComponent
	if ds.db.Preload("Codes").Preload("Configs").Where(&model.GameComponent{UserID: user.ID}).First(&luchador).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchador")

	return &luchador
}

func (ds *DataSource) findLuchadorByIDNoPreload(id uint) *model.GameComponent {
	var luchador model.GameComponent
	if ds.db.First(&luchador, id).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchadorByID")

	return &luchador
}

func (ds *DataSource) updateLuchador(component *model.GameComponent) *model.GameComponent {
	current := ds.findLuchadorByID(component.ID)
	if current == nil {
		return nil
	}

	current.Name = component.Name
	current.Configs = applyConfigChanges(current.Configs, component.Configs)
	current.Codes = component.Codes

	ds.db.Save(current)

	log.WithFields(log.Fields{
		"luchador": current,
	}).Info("after updateLuchador")

	return current
}

func applyConfigChanges(original []model.Config, updated []model.Config) []model.Config {
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

func (ds *DataSource) findActiveMultiplayerMatches() *[]model.Match {

	var matches []model.Match
	ds.db.
		Joins("left join game_definitions on matches.game_definition_id = game_definitions.id").
		Where("game_definitions.type = ?", GAMEDEFINITION_TYPE_MULTIPLAYER).
		Where("time_end < time_start").
		Order("time_start desc").Find(&matches)

	log.WithFields(log.Fields{
		"matches": matches,
	}).Info("findActiveMatches")

	return &matches
}

func (ds *DataSource) findActiveMatches() *[]model.Match {

	var matches []model.Match
	ds.db.Where("time_end < time_start").Order("time_start desc").Find(&matches)

	log.WithFields(log.Fields{
		"matches": matches,
	}).Info("findActiveMatches")

	return &matches
}

// TODO: part of this logic will be moved to the runner to define if there is an match and if the
// is already in that

// func (ds *DataSource) findActiveMatchesByGameDefinitionAndParticipant(gameDefinition *model.GameDefinition, gameComponent *model.GameComponent) *model.Match {

// 	var matches []model.Match
// 	ds.db.Preload("Participants").Where(&model.Match{GameDefinitionID: gameDefinition.ID}).Where("time_end < time_start").Find(&matches)

// 	log.WithFields(log.Fields{
// 		"matches": matches,
// 	}).Info("findActiveMatchesByGameDefinitionAndParticipant")

// 	for _, match := range matches {
// 		for _, participant := range match.Participants {
// 			if participant.ID == gameComponent.ID {
// 				return &match
// 			}
// 		}
// 	}

// 	return nil
// }

func (ds *DataSource) findMaskConfig(id uint) *[]model.Config {

	var component model.GameComponent
	if ds.db.Preload("Configs").Where(&model.GameComponent{ID: id}).First(&component).RecordNotFound() {
		var configs []model.Config
		return &configs
	}

	log.WithFields(log.Fields{
		"luchador": component,
	}).Info("findLuchador")

	log.WithFields(log.Fields{
		"configs": component.Configs,
	}).Info("findMaskConfig")

	return &component.Configs
}

func (ds *DataSource) findMatch(id uint) *model.Match {

	var match model.Match
	ds.db.Preload("Participants").Where(&model.Match{ID: id}).First(&match)

	log.WithFields(log.Fields{
		"id":    id,
		"match": match,
	}).Info("findMatch")

	return &match
}

func (ds *DataSource) findLuchadorByID(luchadorID uint) *model.GameComponent {
	var luchador model.GameComponent
	if ds.db.Preload("Codes").Preload("Configs").Where(&model.GameComponent{ID: luchadorID}).First(&luchador).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchadorByID")

	return &luchador
}

func (ds *DataSource) findLuchadorByName(name string) *model.GameComponent {
	var luchador model.GameComponent
	if ds.db.Where(&model.GameComponent{Name: name}).First(&luchador).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchadorByName")

	return &luchador
}

func (ds *DataSource) findLuchadorByNamePreload(name string) *model.GameComponent {
	var luchador model.GameComponent
	if ds.db.Preload("Codes").Preload("Configs").Where(&model.GameComponent{Name: name}).First(&luchador).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("findLuchadorByNamePreload")

	return &luchador
}

func (ds *DataSource) nameExist(ID uint, name string) bool {
	var luchador model.GameComponent
	result := !ds.db.Where("id <> ? AND name = ?", ID, name).First(&luchador).RecordNotFound()

	log.WithFields(log.Fields{
		"luchador": luchador,
		"result":   result,
	}).Debug("NameExist")

	return result
}

func (ds *DataSource) addMatchParticipant(mp *model.MatchParticipant) *model.MatchParticipant {

	var match *model.Match
	match = ds.findMatch(mp.MatchID)
	if match == nil {
		log.WithFields(log.Fields{
			"matchID": mp.MatchID,
		}).Error("Match not found")
		return nil
	}

	var component *model.GameComponent
	component = ds.findLuchadorByIDNoPreload(mp.LuchadorID)
	if component == nil {
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
			}).Warning("Luchador is already in the match")

			return &(model.MatchParticipant{MatchID: mp.MatchID, LuchadorID: mp.LuchadorID})
		}
	}

	match.Participants = append(match.Participants, *component)
	ds.db.Save(&match)

	matchPartipant := model.MatchParticipant{
		LuchadorID: component.ID,
		MatchID:    match.ID,
	}

	log.WithFields(log.Fields{
		"matchPartipant": matchPartipant,
	}).Info("MatchPartipant created")

	return &matchPartipant
}

func (ds *DataSource) endMatch(match *model.Match) *model.Match {

	ds.db.Model(&match).Update("time_end", match.TimeEnd)

	log.WithFields(log.Fields{
		"match": match,
	}).Info("Match time_end updated")

	return match
}

func (ds *DataSource) findLuchadorConfigsByMatchID(id uint) *[]model.GameComponent {

	match := model.Match{}
	ds.db.First(&match, "id = ?", id)

	var participants []model.GameComponent
	ds.db.Model(&match).Related(&participants, "Participants").Preload("Configs")

	log.WithFields(log.Fields{
		"id":     id,
		"match":  match,
		"result": participants,
	}).Debug("findLuchadorConfigsByMatchID")

	return &participants
}

func (ds *DataSource) getMatchScoresByMatchID(id uint) *[]model.MatchScore {

	result := []model.MatchScore{}
	ds.db.Where(&model.MatchScore{MatchID: id}).Find(&result)

	for _, val := range result {
		log.WithFields(log.Fields{
			"matchId":    val.MatchID,
			"luchadorId": val.LuchadorID,
			"score":      val.Score,
		}).Debug("getMatchScoresByMatchID")
	}

	return &result
}

func (ds *DataSource) addMatchScores(ms *model.ScoreList) *model.ScoreList {

	log.WithFields(log.Fields{
		"action":    "start",
		"scorelist": ms,
	}).Info("addMatchScores")

	var match *model.Match = nil

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

		var component *model.GameComponent
		component = ds.findLuchadorByID(score.LuchadorID)
		if component == nil {
			log.WithFields(log.Fields{
				"luchadorID": score.LuchadorID,
			}).Error("Luchador not found")
			return nil
		}

		log.WithFields(log.Fields{
			"action":   "luchador-found",
			"luchador": component,
		}).Info("addMatchScores")

		score := model.MatchScore{
			LuchadorID: component.ID,
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

func (ds *DataSource) updateGameDefinition(input *model.GameDefinition) *model.GameDefinition {

	gameDefinition := ds.findGameDefinitionByName(input.Name)

	if gameDefinition != nil {

		gameDefinition.Duration = input.Duration
		gameDefinition.MinParticipants = input.MinParticipants
		gameDefinition.MaxParticipants = input.MaxParticipants
		gameDefinition.ArenaWidth = input.ArenaWidth
		gameDefinition.ArenaHeight = input.ArenaHeight
		gameDefinition.BulletSize = input.BulletSize
		gameDefinition.LuchadorSize = input.LuchadorSize
		gameDefinition.Fps = input.Fps
		gameDefinition.BuletSpeed = input.BuletSpeed
		gameDefinition.Label = input.Label
		gameDefinition.Description = input.Description
		gameDefinition.Type = input.Type
		gameDefinition.SortOrder = input.SortOrder
		gameDefinition.RadarAngle = input.RadarAngle
		gameDefinition.RadarRadius = input.RadarRadius
		gameDefinition.PunchAngle = input.PunchAngle
		gameDefinition.Life = input.Life
		gameDefinition.Energy = input.Energy
		gameDefinition.PunchDamage = input.PunchDamage
		gameDefinition.PunchCoolDown = input.PunchCoolDown
		gameDefinition.MoveSpeed = input.MoveSpeed
		gameDefinition.TurnSpeed = input.TurnSpeed
		gameDefinition.TurnGunSpeed = input.TurnGunSpeed
		gameDefinition.RespawnCooldown = input.RespawnCooldown
		gameDefinition.MaxFireCooldown = input.MaxFireCooldown
		gameDefinition.MinFireDamage = input.MinFireDamage
		gameDefinition.MaxFireDamage = input.MaxFireDamage
		gameDefinition.MinFireAmount = input.MinFireAmount
		gameDefinition.MaxFireAmount = input.MaxFireAmount
		gameDefinition.RestoreEnergyperSecond = input.RestoreEnergyperSecond
		gameDefinition.RecycledLuchadorEnergyRestore = input.RecycledLuchadorEnergyRestore
		gameDefinition.IncreaseSpeedEnergyCost = input.IncreaseSpeedEnergyCost
		gameDefinition.IncreaseSpeedPercentage = input.IncreaseSpeedPercentage
		gameDefinition.FireEnergyCost = input.FireEnergyCost
		gameDefinition.RespawnX = input.RespawnX
		gameDefinition.RespawnY = input.RespawnY
		gameDefinition.RespawnAngle = input.RespawnAngle
		gameDefinition.RespawnGunAngle = input.RespawnGunAngle

		dbc := ds.db.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "save",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		for n, gc := range input.GameComponents {
			component := ds.findLuchadorByNamePreload(gc.Name)

			log.WithFields(log.Fields{
				"gc.Name": gc.Name,
				"gc.ID":   gc.ID,
			}).Debug("searching gamedefinition")

			if component != nil {
				input.GameComponents[n] = *(ds.updateLuchador(component))
			}
		}

		ds.db.Model(gameDefinition).Association("GameComponents").Replace(input.GameComponents)
		dbc = ds.db.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "GameComponents",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		ds.db.Model(gameDefinition).Association("SceneComponents").Replace(input.SceneComponents)
		dbc = ds.db.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "SceneComponents",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		ds.db.Model(gameDefinition).Association("Codes").Replace(input.Codes)
		dbc = ds.db.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "Codes",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		ds.db.Model(gameDefinition).Association("LuchadorSuggestedCodes").Replace(input.LuchadorSuggestedCodes)
		dbc = ds.db.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "LuchadorSuggestedCodes",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		log.WithFields(log.Fields{
			"gameDefinition": gameDefinition,
		}).Info("updateGameDefinition")

		return gameDefinition
	}

	// not found
	return nil

}

func (ds *DataSource) createGameDefinition(g *model.GameDefinition) *model.GameDefinition {

	gameDefinition := model.GameDefinition{}
	copier.Copy(&gameDefinition, &g)
	for n := range g.GameComponents {
		g.GameComponents[n].Configs = randomConfig()
	}

	ds.db.Create(&gameDefinition)

	log.WithFields(log.Fields{
		"gameDefinition": gameDefinition,
	}).Info("createGameDefinition")

	return &gameDefinition
}

func (ds *DataSource) findGameDefinition(id uint) *model.GameDefinition {
	var gameDefinition model.GameDefinition

	if ds.db.
		Preload("GameComponents").
		Preload("GameComponents.Codes").
		Preload("GameComponents.Configs").
		Preload("SceneComponents").
		Preload("SceneComponents.Codes").
		Preload("Codes").
		Preload("LuchadorSuggestedCodes").
		Where(&model.GameDefinition{ID: id}).
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
	}).Debug("findGameDefinition before array checks")

	resetGameDefinitionArrays(&gameDefinition)

	log.WithFields(log.Fields{
		"ID":             id,
		"gameDefinition": gameDefinition,
	}).Info("findGameDefinition")

	return &gameDefinition
}

func (ds *DataSource) findGameDefinitionByName(name string) *model.GameDefinition {
	var gameDefinition model.GameDefinition
	var filter model.GameDefinition
	filter.Name = name

	if ds.db.
		Preload("GameComponents").
		Preload("GameComponents.Codes").
		Preload("GameComponents.Configs").
		Preload("SceneComponents").
		Preload("SceneComponents.Codes").
		Preload("Codes").
		Preload("LuchadorSuggestedCodes").
		Where(&filter).
		First(&gameDefinition).
		RecordNotFound() {

		log.WithFields(log.Fields{
			"Name":           name,
			"gamedefinition": gameDefinition,
		}).Error("findGameDefinitionByName not found")

		return nil
	}

	log.WithFields(log.Fields{
		"Name":           name,
		"gameDefinition": gameDefinition,
	}).Debug("findGameDefinitionByName before array checks")

	resetGameDefinitionArrays(&gameDefinition)

	log.WithFields(log.Fields{
		"Name":           name,
		"gameDefinition": gameDefinition,
	}).Debug("findGameDefinitionByName")

	return &gameDefinition
}

func (ds *DataSource) findAllGameDefinition() *[]model.GameDefinition {
	var gameDefinitions []model.GameDefinition

	ds.db.
		Preload("GameComponents").
		Preload("GameComponents.Codes").
		Preload("GameComponents.Configs").
		Preload("SceneComponents").
		Preload("SceneComponents.Codes").
		Preload("Codes").
		Preload("LuchadorSuggestedCodes").
		Order("sort_order").
		Find(&gameDefinitions)

	log.WithFields(log.Fields{
		"gameDefinitions": gameDefinitions,
	}).Debug("findTutorialGameDefinition before array checks")

	for i := range gameDefinitions {
		resetGameDefinitionArrays(&gameDefinitions[i])
	}

	log.WithFields(log.Fields{
		"gameDefinitions": gameDefinitions,
	}).Debug("findAllGameDefinition")

	return &gameDefinitions
}

func (ds *DataSource) findTutorialGameDefinition() *[]model.GameDefinition {
	var gameDefinitions []model.GameDefinition

	var filter model.GameDefinition
	filter.Type = GAMEDEFINITION_TYPE_TUTORIAL

	ds.db.
		Preload("GameComponents").
		Preload("GameComponents.Codes").
		Preload("GameComponents.Configs").
		Preload("SceneComponents").
		Preload("SceneComponents.Codes").
		Preload("Codes").
		Preload("LuchadorSuggestedCodes").
		Where(&filter).
		Order("sort_order").
		Find(&gameDefinitions)

	log.WithFields(log.Fields{
		"gameDefinitions": gameDefinitions,
	}).Debug("findTutorialGameDefinition before array checks")

	for i := range gameDefinitions {
		resetGameDefinitionArrays(&gameDefinitions[i])
	}

	log.WithFields(log.Fields{
		"gameDefinitions": gameDefinitions,
	}).Debug("findTutorialGameDefinition")

	return &gameDefinitions
}

func resetGameDefinitionArrays(gameDefinition *model.GameDefinition) {
	if gameDefinition.GameComponents == nil {
		gameDefinition.GameComponents = make([]model.GameComponent, 0)
	}

	if gameDefinition.SceneComponents == nil {
		gameDefinition.SceneComponents = make([]model.SceneComponent, 0)
	}

	if gameDefinition.Codes == nil {
		gameDefinition.Codes = make([]model.Code, 0)
	}

	if gameDefinition.LuchadorSuggestedCodes == nil {
		gameDefinition.LuchadorSuggestedCodes = make([]model.Code, 0)
	}
}

func (ds *DataSource) addMatchMetric(m *model.MatchMetric) *model.MatchMetric {

	metric := model.MatchMetric{}
	copier.Copy(&metric, &m)
	ds.db.Create(&metric)

	log.WithFields(log.Fields{
		"metric": metric,
	}).Debug("addMatchMetric")

	return &metric
}

// TODO: remove this
var accessCodeCounter int64

func (ds *DataSource) addClassroom(c *model.Classroom) *model.Classroom {

	now := fmt.Sprintf("%X", time.Now().Unix()+accessCodeCounter)
	accessCodeCounter = accessCodeCounter + 1

	classroom := model.Classroom{
		Name:       c.Name,
		OwnerID:    c.OwnerID,
		AccessCode: now,
	}

	log.WithFields(log.Fields{
		"classroom": classroom,
	}).Debug("addClassroom")

	ds.db.Create(&classroom)
	classroom.Students = make([]model.Student, 0)

	log.WithFields(log.Fields{
		"classroom": classroom,
	}).Debug("after addClassroom")

	return &classroom
}

func (ds *DataSource) findAllClassroom(user *model.User) *[]model.Classroom {
	var result []model.Classroom

	ds.db.
		Preload("Students").
		Where(&model.Classroom{OwnerID: user.ID}).
		Order("name").
		Find(&result)

	log.WithFields(log.Fields{
		"classrooms": result,
	}).Debug("findAllClassroom")

	log.WithFields(log.Fields{
		"classrooms": result,
	}).Debug("findAllClassroom")

	return &result
}

func (ds *DataSource) joinClassroom(user *model.User, accessCode string) *model.Classroom {
	var result model.Classroom
	student := model.Student{UserID: user.ID}

	if ds.db.Preload("Students").
		Where(&model.Classroom{AccessCode: accessCode}).
		First(&result).
		RecordNotFound() {

		log.WithFields(log.Fields{
			"accessCode": accessCode,
		}).Info("classroom not found")

		return nil
	}

	if ds.db.
		Where(&student).
		First(&student).
		RecordNotFound() {

		log.WithFields(log.Fields{
			"userID": user.ID,
		}).Info("student not found will create")

		ds.db.Create(&student)
	}

	result.Students = append(result.Students, student)
	ds.db.Save(&result)

	log.WithFields(log.Fields{
		"classroom": result,
	}).Debug("joinClassroom")

	return &result
}

func (ds *DataSource) findPublicAvailableMatch() *[]model.AvailableMatch {
	return ds.findAvailableMatchByClassroomID(0)
}

func (ds *DataSource) findAvailableMatchByClassroomID(id uint) *[]model.AvailableMatch {

	var result []model.AvailableMatch

	ds.db.
		Where("classroom_id == ?", id).
		Find(&result)

	log.WithFields(log.Fields{
		"availableMatch": result,
	}).Debug("findAvailableMatchByClassroomID")

	return &result
}
