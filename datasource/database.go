package datasource

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
	DB     *gorm.DB
	secret string
}

// BuildMysqlConfig creates a DBconfig for Mysql based on environment variables
func BuildMysqlConfig() *DBconfig {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
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

// BuildSQLLiteConfig creates a DBconfig for Mysql based on environment variables
func BuildSQLLiteConfig(fileName string) *DBconfig {
	return &DBconfig{
		dialect: "sqlite3",
		args:    fileName}
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
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Session{})
	DB.AutoMigrate(&model.UserSetting{})
	DB.AutoMigrate(&model.UserLevel{})
	DB.AutoMigrate(&model.Match{})
	DB.AutoMigrate(&model.Code{})
	DB.AutoMigrate(&model.CodeHistory{})

	DB.AutoMigrate(&model.Config{})
	DB.AutoMigrate(&model.MatchScore{})
	DB.AutoMigrate(&model.SceneComponent{})
	DB.AutoMigrate(&model.GameComponent{})
	DB.AutoMigrate(&model.GameDefinition{})
	DB.AutoMigrate(&model.Classroom{})
	DB.AutoMigrate(&model.Student{})
	DB.AutoMigrate(&model.AvailableMatch{})

	DB.AutoMigrate(&model.LearningObjective{})
	DB.AutoMigrate(&model.Skill{})
	// DB.AutoMigrate(&model.GradingSystem{})
	DB.AutoMigrate(&model.Grade{})
	DB.AutoMigrate(&model.Activity{})

	secret := os.Getenv("API_SECRET")

	return &DataSource{DB: DB, config: config, secret: secret}
}

// KeepAlive sends ticks to the DB to keep the connection alive
func (ds *DataSource) KeepAlive() {
	log.Debug("Keep connection alive")
	for range time.Tick(time.Minute) {
		ds.DB.DB().Ping()
		log.Debug("Keep connection alive")
	}
}

func (ds *DataSource) findUserByEmail(email string) *model.User {
	var user model.User

	if ds.DB.Where("email = ?", email).First(&user).RecordNotFound() {
		return nil
	}
	return &user
}

func (ds *DataSource) findUserBySession(UUID string) *model.User {
	var session model.Session
	var user model.User

	if ds.DB.Where("UUID = ?", UUID).First(&session).RecordNotFound() {
		return nil
	}

	if ds.DB.Where("ID = ?", session.UserID).First(&user).RecordNotFound() {
		return nil
	}

	return &user
}

// FindUserByID definition
func (ds *DataSource) FindUserByID(id uint) *model.User {
	var user model.User
	if ds.DB.Where(&model.User{ID: id}).First(&user).RecordNotFound() {
		return nil
	}
	return &user
}

// Create if doesnt exist
func (ds *DataSource) FindUserSettingByUser(user *model.User) *model.UserSetting {
	var settings model.UserSetting
	ds.DB.Where(&model.UserSetting{UserID: user.ID}).FirstOrCreate(&settings)
	return &settings
}

// Create if doesnt exist
func (ds *DataSource) FindUserLevelByUser(user *model.User) *model.UserLevel {
	var level model.UserLevel
	ds.DB.Where(&model.UserLevel{UserID: user.ID}).FirstOrCreate(&level)
	return &level
}

func (ds *DataSource) UpdateUserLevel(level *model.UserLevel) *model.UserLevel {
	var current model.UserLevel
	if ds.DB.First(&current, level.ID).RecordNotFound() {
		return nil
	}

	current.Level = level.Level
	ds.DB.Save(&current)

	log.WithFields(log.Fields{
		"user level": current,
	}).Error("User Level updated")

	return &current
}

func (ds *DataSource) UpdateUserSetting(settings *model.UserSetting) *model.UserSetting {
	var current model.UserSetting
	if ds.DB.First(&current, settings.ID).RecordNotFound() {
		return nil
	}

	current.VisitedMainPage = settings.VisitedMainPage
	current.VisitedMaskPage = settings.VisitedMaskPage
	current.PlayedTutorial = settings.PlayedTutorial

	ds.DB.Save(&current)

	log.WithFields(log.Fields{
		"settings": current,
	}).Error("User Setting updated")

	return &current
}

// CreateUser definition
func (ds *DataSource) CreateUser(name string) *model.User {
	user := model.User{Username: name}
	ds.DB.Where(&model.User{Username: name}).FirstOrCreate(&user)

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

func (ds *DataSource) CreateLuchador(l *model.GameComponent) *model.GameComponent {
	luchador := model.GameComponent{
		UserID:  l.UserID,
		Name:    l.Name,
		Codes:   l.Codes,
		Configs: l.Configs,
	}

	ds.DB.Create(&luchador)

	log.WithFields(log.Fields{
		"luchador": model.LogGameComponent(&luchador),
	}).Info("Luchador created")

	return &luchador
}

func (ds *DataSource) FindLuchador(user *model.User) *model.GameComponent {
	var luchador model.GameComponent
	if ds.DB.Preload("Codes").Preload("Configs").Where(&model.GameComponent{UserID: user.ID}).First(&luchador).RecordNotFound() {
		return nil
	}

	luchador.Codes = removeDuplicates(luchador.Codes)

	log.WithFields(log.Fields{
		"luchador": model.LogGameComponent(&luchador),
	}).Info("FindLuchador")

	return &luchador
}

func (ds *DataSource) FindLuchadorByIDNoPreload(id uint) *model.GameComponent {
	var luchador model.GameComponent
	if ds.DB.First(&luchador, id).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("FindLuchadorByID")

	return &luchador
}

// UpdateLuchador keeping unique config an code
func (ds *DataSource) UpdateLuchador(component *model.GameComponent) *model.GameComponent {
	current := ds.FindLuchadorByID(component.ID)
	if current == nil {
		return nil
	}

	current.Name = component.Name
	current.Configs = applyConfigChanges(current.Configs, component.Configs)

	current.Codes = applyCodeChanges(current.Codes, component.Codes)
	ds.DB.Save(current)

	log.WithFields(log.Fields{
		"luchador": current,
	}).Warning("after updateLuchador")

	return current
}

func applyCodeChanges(current []model.Code, updated []model.Code) []model.Code {
	var found bool

	for _, newCode := range updated {
		found = false

		//search by event+gameDefinition
		for i, currentCode := range current {
			if currentCode.Event == newCode.Event && currentCode.GameDefinitionID == newCode.GameDefinitionID {
				current[i].Script = newCode.Script
				found = true
				break
			}
		}

		if !found {
			current = append(current, newCode)
		}
	}

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

// func (ds *DataSource) ActiveMatchesSQL() string {

// 	endBeforeStart := "time_end < time_start"
// 	start := "UNIX_TIMESTAMP(time_start)"
// 	duration := "(game_definitions.duration/1000)"
// 	now := "UNIX_TIMESTAMP(NOW())"

// 	sql := fmt.Sprintf("%v and %v + %v > %v",
// 		endBeforeStart, start, duration, now)

// 	return sql
// }

func (ds *DataSource) FindActiveMultiplayerMatches() *[]model.Match {

	matches := ds.FindActiveMatches("game_definitions.type = ?", model.GAMEDEFINITION_TYPE_MULTIPLAYER)

	log.WithFields(log.Fields{
		"matches": model.LogMatches(matches),
	}).Info("findActiveMatches")

	return matches
}

// FindActiveMatches definition
func (ds *DataSource) FindActiveMatches(query interface{}, args ...interface{}) *[]model.Match {

	var matches []model.Match
	ds.DB.
		Joins("left join game_definitions on matches.game_definition_id = game_definitions.id").
		Preload("GameDefinition").
		Preload("Participants").
		Where("time_end < time_start").
		Where(query, args).
		Order("time_start desc").
		Find(&matches)

	log.WithFields(log.Fields{
		"matches": matches,
	}).Warn("findActiveMatches")

	result := make([]model.Match, 0)

	// auto remove matches where the duration is greater than the current time
	for _, match := range matches {
		// only if the match gamedefinition has duration
		if match.GameDefinition.Duration > 0 {
			duration := time.Duration(match.GameDefinition.Duration) * time.Millisecond
			startPlusDuration := match.TimeStart.Add(duration)
			now := time.Now()

			log.WithFields(log.Fields{
				"duration":          duration,
				"startPlusDuration": startPlusDuration,
				"now":               now,
				"isAfter":           startPlusDuration.After(now),
			}).Debug("findActiveMatches/time")

			if startPlusDuration.After(now) {
				result = append(result, match)
			}
		} else {
			result = append(result, match)
		}
	}

	return &result
}

// func (ds *DataSource) FindActiveMatchesByGameDefinitionAndParticipant(gameDefinition *model.GameDefinition, gameComponent *model.GameComponent) *model.Match {

// 	matches := ds.FindActiveMatches(&model.Match{GameDefinitionID: gameDefinition.ID})

// 	// var matches []model.Match
// 	// ds.DB.Preload("Participants").
// 	// 	Joins("left join game_definitions on matches.game_definition_id = game_definitions.id").
// 	// 	Where(&model.Match{GameDefinitionID: gameDefinition.ID}).
// 	// 	Where(ds.ActiveMatchesSQL()).
// 	// 	Find(&matches)

// 	log.WithFields(log.Fields{
// 		"matches": matches,
// 	}).Info("findActiveMatchesByGameDefinitionAndParticipant")

// 	for _, match := range *matches {
// 		for _, participant := range match.Participants {
// 			if participant.ID == gameComponent.ID {
// 				return &match
// 			}
// 		}
// 	}

// 	return nil
// }

func (ds *DataSource) FindMaskConfig(id uint) *[]model.Config {

	var component model.GameComponent
	if ds.DB.Preload("Configs").Where(&model.GameComponent{ID: id}).First(&component).RecordNotFound() {
		var configs []model.Config
		return &configs
	}

	log.WithFields(log.Fields{
		"luchador": component,
	}).Info("FindLuchador")

	log.WithFields(log.Fields{
		"configs": component.Configs,
	}).Info("findMaskConfig")

	return &component.Configs
}

func (ds *DataSource) FindMatch(id uint) *model.Match {

	var match model.Match
	ds.DB.Preload("Participants").Where(&model.Match{ID: id}).First(&match)

	log.WithFields(log.Fields{
		"id":    id,
		"match": match,
	}).Info("FindMatch")

	return &match
}

func (ds *DataSource) FindLuchadorByID(luchadorID uint) *model.GameComponent {
	var luchador model.GameComponent
	if ds.DB.Preload("Codes").Preload("Configs").Where(&model.GameComponent{ID: luchadorID}).First(&luchador).RecordNotFound() {
		return nil
	}

	luchador.Codes = removeDuplicates(luchador.Codes)

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("FindLuchadorByID")

	return &luchador
}

func (ds *DataSource) FindLuchadorByName(name string) *model.GameComponent {
	var luchador model.GameComponent
	if ds.DB.Where(&model.GameComponent{Name: name}).First(&luchador).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("FindLuchadorByName")

	return &luchador
}

func (ds *DataSource) FindLuchadorByNamePreload(name string) *model.GameComponent {
	var luchador model.GameComponent
	if ds.DB.Preload("Codes").Preload("Configs").Where(&model.GameComponent{Name: name}).First(&luchador).RecordNotFound() {
		return nil
	}

	luchador.Codes = removeDuplicates(luchador.Codes)

	log.WithFields(log.Fields{
		"luchador": luchador,
	}).Info("FindLuchadorByNamePreload")

	return &luchador
}

func removeDuplicates(current []model.Code) []model.Code {
	var found bool
	result := make([]model.Code, 0)

	for _, code := range current {
		found = false

		//search by event+gameDefinition
		for i, newCode := range result {
			if newCode.Event == code.Event && newCode.GameDefinitionID == code.GameDefinitionID {
				if code.ID > newCode.ID {
					result[i] = code
				}

				found = true
				break
			}
		}

		if !found {
			result = append(result, code)
		}
	}

	return result
}

func (ds *DataSource) NameExist(ID uint, name string) bool {
	var luchador model.GameComponent
	result := !ds.DB.Where("id <> ? AND name = ?", ID, name).First(&luchador).RecordNotFound()

	log.WithFields(log.Fields{
		"luchador": luchador,
		"result":   result,
	}).Debug("NameExist")

	return result
}

func (ds *DataSource) AddMatchParticipant(mp *model.MatchParticipant) *model.MatchParticipant {

	var match *model.Match
	match = ds.FindMatch(mp.MatchID)
	if match == nil {
		log.WithFields(log.Fields{
			"matchID": mp.MatchID,
		}).Error("Match not found")
		return nil
	}

	var component *model.GameComponent
	component = ds.FindLuchadorByIDNoPreload(mp.LuchadorID)
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
	ds.DB.Save(&match)

	matchPartipant := model.MatchParticipant{
		LuchadorID: component.ID,
		MatchID:    match.ID,
	}

	log.WithFields(log.Fields{
		"matchPartipant": matchPartipant,
	}).Info("MatchPartipant created")

	return &matchPartipant
}

func (ds *DataSource) EndMatch(match *model.Match) *model.Match {

	ds.DB.Model(&match).Update("time_end", match.TimeEnd)

	log.WithFields(log.Fields{
		"match": model.LogMatch(match),
	}).Info("Match time_end updated")

	return match
}

func (ds *DataSource) FindLuchadorConfigsByMatchID(id uint) *[]model.GameComponent {

	match := model.Match{}
	ds.DB.First(&match, "id = ?", id)

	var participants []model.GameComponent
	ds.DB.Model(&match).Related(&participants, "Participants").Preload("Configs")

	log.WithFields(log.Fields{
		"id":     id,
		"match":  match,
		"result": participants,
	}).Debug("FindLuchadorConfigsByMatchID")

	return &participants
}

func (ds *DataSource) GetMatchScoresByMatchID(id uint) *[]model.MatchScore {

	result := []model.MatchScore{}
	ds.DB.Where(&model.MatchScore{MatchID: id}).Find(&result)

	for _, val := range result {
		log.WithFields(log.Fields{
			"matchId":    val.MatchID,
			"luchadorId": val.LuchadorID,
			"score":      val.Score,
		}).Debug("getMatchScoresByMatchID")
	}

	return &result
}

func (ds *DataSource) AddMatchScores(ms *model.ScoreList) *model.ScoreList {

	log.WithFields(log.Fields{
		"action":    "start",
		"scorelist": ms,
	}).Info("addMatchScores")

	var match *model.Match = nil

	for _, score := range ms.Scores {

		if match == nil {
			match = ds.FindMatch(score.MatchID)
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
		component = ds.FindLuchadorByID(score.LuchadorID)
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

		ds.DB.Create(&score)

		log.WithFields(log.Fields{
			"action": "after-save",
			"score":  score,
		}).Info("addMatchScores")
	}

	return ms
}

func (ds *DataSource) UpdateGameDefinition(input *model.GameDefinition) *model.GameDefinition {

	gameDefinition := ds.FindGameDefinitionByName(input.Name)

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

		dbc := ds.DB.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "save",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		for n, gc := range input.GameComponents {
			component := ds.FindLuchadorByNamePreload(gc.Name)

			log.WithFields(log.Fields{
				"gc.Name": gc.Name,
				"gc.ID":   gc.ID,
			}).Debug("searching gamedefinition")

			if component != nil {
				input.GameComponents[n] = *(ds.UpdateLuchador(component))
			}
		}

		ds.DB.Model(gameDefinition).Association("GameComponents").Replace(input.GameComponents)
		dbc = ds.DB.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "GameComponents",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		ds.DB.Model(gameDefinition).Association("SceneComponents").Replace(input.SceneComponents)
		dbc = ds.DB.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "SceneComponents",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		ds.DB.Model(gameDefinition).Association("Codes").Replace(input.Codes)
		dbc = ds.DB.Save(gameDefinition)
		if dbc.Error != nil {
			log.WithFields(log.Fields{
				"error":               dbc.Error,
				"gameDefinition.Name": gameDefinition.Name,
				"step":                "Codes",
			}).Error("Error updating updateGameDefinition")

			return nil
		}

		ds.DB.Model(gameDefinition).Association("LuchadorSuggestedCodes").Replace(input.LuchadorSuggestedCodes)
		dbc = ds.DB.Save(gameDefinition)
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

func (ds *DataSource) CreateGameDefinition(g *model.GameDefinition) *model.GameDefinition {

	gameDefinition := model.GameDefinition{}
	copier.Copy(&gameDefinition, &g)
	for n := range g.GameComponents {
		g.GameComponents[n].Configs = model.RandomConfig()
	}

	ds.DB.Create(&gameDefinition)

	log.WithFields(log.Fields{
		"gameDefinition": gameDefinition,
	}).Info("createGameDefinition")

	return &gameDefinition
}

func (ds *DataSource) FindGameDefinition(id uint) *model.GameDefinition {
	var gameDefinition model.GameDefinition

	if ds.DB.
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

// FindGameDefinitionByName definition
func (ds *DataSource) FindGameDefinitionByName(name string) *model.GameDefinition {
	var gameDefinition model.GameDefinition

	if ds.DB.
		Preload("GameComponents").
		Preload("GameComponents.Codes").
		Preload("GameComponents.Configs").
		Preload("SceneComponents").
		Preload("SceneComponents.Codes").
		Preload("Codes").
		Preload("LuchadorSuggestedCodes").
		Where(&model.GameDefinition{Name: name}).
		First(&gameDefinition).
		RecordNotFound() {

		log.WithFields(log.Fields{
			"Name":           name,
			"gamedefinition": gameDefinition,
		}).Error("FindGameDefinitionByName not found")

		return nil
	}

	log.WithFields(log.Fields{
		"Name":           name,
		"gameDefinition": gameDefinition,
	}).Debug("FindGameDefinitionByName before array checks")

	resetGameDefinitionArrays(&gameDefinition)

	log.WithFields(log.Fields{
		"Name":           name,
		"gameDefinition": gameDefinition,
	}).Debug("FindGameDefinitionByName")

	return &gameDefinition
}

func (ds *DataSource) FindAllGameDefinition() *[]model.GameDefinition {
	var gameDefinitions []model.GameDefinition

	ds.DB.
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

func (ds *DataSource) FindTutorialGameDefinition() *[]model.GameDefinition {
	var gameDefinitions []model.GameDefinition

	ds.DB.
		Preload("GameComponents").
		Preload("GameComponents.Codes").
		Preload("GameComponents.Configs").
		Preload("SceneComponents").
		Preload("SceneComponents.Codes").
		Preload("Codes").
		Preload("LuchadorSuggestedCodes").
		Where(&model.GameDefinition{Type: model.GAMEDEFINITION_TYPE_TUTORIAL}).
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

// CreateAvailableMatchIfDontExist definition
func (ds *DataSource) CreateAvailableMatchIfDontExist(gameDefinitionID uint, name string) *model.AvailableMatch {
	am := model.AvailableMatch{GameDefinitionID: gameDefinitionID, Name: name}
	ds.DB.Where(&am).FirstOrCreate(&am)

	log.WithFields(log.Fields{
		"gameDefinitionID": gameDefinitionID,
		"AvailableMatch":   am,
	}).Debug("CreateAvailableMatchIfDontExist")

	return &am
}
