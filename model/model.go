package model

import (
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// User definition
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	Username  string     `json:"username"`
}

// UserDetails definition
type UserDetails struct {
	User       *User       `json:"user"`
	Classrooms []Classroom `json:"classrooms"`
	Roles      []string    `json:"roles"`
	Settings   UserSetting `json:"settings"`
	Level      UserLevel   `json:"level"`
}

// Session definition
type Session struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	UUID      string     `json:"UUID"`
	UserID    uint       `json:"userID"`
}

// UserSetting definition
type UserSetting struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `json:"-"`
	UpdatedAt       time.Time  `json:"-"`
	DeletedAt       *time.Time `json:"-" faker:"-"`
	UserID          uint       `json:"userID"`
	VisitedMainPage bool       `json:"visitedMainPage"`
	VisitedMaskPage bool       `json:"visitedMaskPage"`
	PlayedTutorial  bool       `json:"playedTutorial"`
}

// UserLevel definition
type UserLevel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	UserID    uint       `json:"userID"`
	Level     uint       `json:"level"`
}

// ActiveMatch definition, describes the result of findActiveMatches
// mixing Tutorial and PVP matches
type ActiveMatch struct {
	MatchID     uint      `json:"matchID"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	SortOrder   uint      `json:"sortOrder"`
	Duration    uint64    `json:"duration"`
	TimeStart   time.Time `json:"timeStart"`
}

// GameDefinition definition
type GameDefinition struct {
	ID                            uint             `gorm:"primary_key" json:"id" faker:"-"`
	CreatedAt                     time.Time        `json:"-"`
	UpdatedAt                     time.Time        `json:"-"`
	DeletedAt                     *time.Time       `json:"-" faker:"-"`
	Duration                      uint64           `json:"duration"`
	MinParticipants               uint             `json:"minParticipants"`
	MaxParticipants               uint             `json:"maxParticipants"`
	ArenaWidth                    uint             `json:"arenaWidth"`
	ArenaHeight                   uint             `json:"arenaHeight"`
	BulletSize                    uint             `json:"bulletSize"`
	LuchadorSize                  uint             `json:"luchadorSize"`
	Fps                           uint             `json:"fps"`
	BuletSpeed                    uint             `json:"buletSpeed"`
	Name                          string           `gorm:"not null;unique_index" json:"name"`
	Label                         string           `json:"label"`
	Description                   string           `json:"description"`
	Type                          string           `json:"type"`
	SortOrder                     uint             `json:"sortOrder"`
	RadarAngle                    uint             `json:"radarAngle"`
	RadarRadius                   uint             `json:"radarRadius"`
	PunchAngle                    uint             `json:"punchAngle"`
	Life                          uint             `json:"life"`
	Energy                        uint             `json:"energy"`
	PunchDamage                   uint             `json:"punchDamage"`
	PunchCoolDown                 uint             `json:"punchCoolDown"`
	MoveSpeed                     uint             `json:"moveSpeed"`
	TurnSpeed                     uint             `json:"turnSpeed"`
	TurnGunSpeed                  uint             `json:"turnGunSpeed"`
	RespawnCooldown               uint             `json:"respawnCooldown"`
	MaxFireCooldown               uint             `json:"maxFireCooldown"`
	MinFireDamage                 uint             `json:"minFireDamage"`
	MaxFireDamage                 uint             `json:"maxFireDamage"`
	MinFireAmount                 uint             `json:"minFireAmount"`
	MaxFireAmount                 uint             `json:"maxFireAmount"`
	RestoreEnergyperSecond        uint             `json:"restoreEnergyperSecond"`
	RecycledLuchadorEnergyRestore uint             `json:"recycledLuchadorEnergyRestore"`
	IncreaseSpeedEnergyCost       uint             `json:"increaseSpeedEnergyCost"`
	IncreaseSpeedPercentage       uint             `json:"increaseSpeedPercentage"`
	FireEnergyCost                uint             `json:"fireEnergyCost"`
	RespawnX                      uint             `json:"respawnX"`
	RespawnY                      uint             `json:"respawnY"`
	RespawnAngle                  uint             `json:"respawnAngle"`
	RespawnGunAngle               uint             `json:"respawnGunAngle"`
	MinLevel                      uint             `json:"minLevel"`
	MaxLevel                      uint             `json:"maxLevel"`
	UnblockLevel                  uint             `json:"unblockLevel"`
	TeamDefinition                TeamDefinition   `json:"teamDefinition"`
	GameComponents                []GameComponent  `json:"gameComponents"`
	SceneComponents               []SceneComponent `json:"sceneComponents"`
	Codes                         []Code           `gorm:"many2many:gamedefinition_codes" json:"codes"`
	LuchadorSuggestedCodes        []Code           `gorm:"many2many:gamedefinition_suggestedcodes" json:"suggestedCodes"`
}

type TeamDefinition struct {
	ID               uint       `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	GameDefinitionID uint       `json:"gameDefinition,omitempty" faker:"-"`
	FriendlyFire     bool       `json:"friendlyFire"`
	Teams            []Team     `json:"teams"`
}

type Team struct {
	ID               uint       `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	TeamDefinitionID uint       `json:"teamDefinition,omitempty" faker:"-"`
	Name             string     `json:"name"`
	Color            string     `json:"color"`
	MinParticipants  uint       `json:"minParticipants"`
	MaxParticipants  uint       `json:"maxParticipants"`
}

// Match definition
type Match struct {
	ID                 uint            `gorm:"primary_key" json:"id"`
	CreatedAt          time.Time       `json:"-"`
	UpdatedAt          time.Time       `json:"-"`
	DeletedAt          *time.Time      `json:"-" faker:"-"`
	AvailableMatchID   uint            `gorm:"unique_index:idx_available_match_state" json:"availableMatchID"`
	State              string          `gorm:"unique_index:idx_available_match_state;default:'created'" json:"state"`
	TimeStart          time.Time       `json:"timeStart"`
	TimeEnd            time.Time       `json:"timeEnd"`
	LastTimeAlive      time.Time       `json:"lastTimeAlive"`
	GameDefinitionID   uint            `json:"gameDefinitionID"`
	GameDefinition     GameDefinition  `json:"gameDefinition"`
	GameDefinitionData string          `gorm:"size:125000" json:"-"`
	Participants       []GameComponent `gorm:"many2many:match_participants" json:"participants"`
}

// SceneComponent definition
type SceneComponent struct {
	ID               uint       `gorm:"primary_key" json:"id" faker:"-"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	GameDefinitionID uint       `json:"gameDefinition,omitempty" faker:"-"`
	X                uint       `json:"x"`
	Y                uint       `json:"y"`
	Width            uint       `json:"width"`
	Height           uint       `json:"height"`
	Rotation         uint       `json:"rotation"`
	Life             uint       `json:"life"`
	Respawn          bool       `json:"respawn"`
	Colider          bool       `json:"colider"`
	ShowInRadar      bool       `json:"showInRadar"`
	BlockMovement    bool       `json:"blockMovement"`
	Type             string     `json:"type"`
	Color            string     `json:"color"`
	Alpha            float32    `json:"alpha"`
	Codes            []Code     `gorm:"many2many:scenecomponent_codes" json:"codes"`
}

// GameComponent definition
type GameComponent struct {
	ID               uint       `gorm:"primary_key" json:"id" faker:"-"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	GameDefinitionID uint       `json:"gameDefinition,omitempty" faker:"-"`
	Name             string     `gorm:"not null;unique_index" json:"name"`
	UserID           uint       `json:"userID,omitempty"`
	Codes            []Code     `gorm:"many2many:gamecomponent_codes" json:"codes"`
	Configs          []Config   `gorm:"many2many:gamecomponent_configs" json:"configs"`
	IsNPC            bool       `json:"isNPC"`
	X                uint       `json:"x"`
	Y                uint       `json:"y"`
	Life             uint       `json:"life"`
	Angle            uint       `json:"angle"`
	GunAngle         uint       `json:"gunAngle"`
}

// Code definition
type Code struct {
	ID               uint       `gorm:"primary_key" json:"id,omitempty" faker:"-"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	Event            string     `json:"event"`
	Script           string     `json:"script"`
	Version          uint       `json:"version"`
	Exception        string     `json:"exception"`
	GameDefinitionID uint       `json:"gameDefinition" faker:"-"`
}

// BeforeSave Code hook, used to version Code
func (c *Code) BeforeSave() (err error) {
	if c.Version == 0 {
		c.Version = 1
	} else {
		c.Version = c.Version + 1
	}

	return
}

// AfterCreate Code hook
func (c *Code) AfterCreate(scope *gorm.Scope) (err error) {
	version := CodeHistory{Script: c.Script, Version: c.Version, CodeID: c.ID}

	log.WithFields(log.Fields{
		"version": version,
	}).Info("Code afterCreate")

	scope.DB().Model(&version).Create(&version)
	return
}

// AfterUpdate Code hook
func (c *Code) AfterUpdate(scope *gorm.Scope) (err error) {
	version := CodeHistory{Script: c.Script, Version: c.Version, CodeID: c.ID}

	log.WithFields(log.Fields{
		"version": version,
	}).Info("Code afterUpdate")

	scope.DB().Model(&version).Create(&version)
	return
}

// CodeHistory definition
type CodeHistory struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty" faker:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	Script    string     `json:"script"`
	Version   uint       `json:"version"`
	CodeID    uint       `json:"codeID"`
	Code      *Code      `json:"code,omitempty"`
}

// Config definition
type Config struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty" faker:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	Key       string     `json:"key"`
	Value     string     `json:"value"`
}

// JoinMatch definition
type JoinMatch struct {
	MatchID    uint `json:"matchID"`
	LuchadorID uint `json:"luchadorID"`
}

// FindLuchadorWithGamedefinition definition
type FindLuchadorWithGamedefinition struct {
	GameDefinitionID uint `json:"gameDefinitionID"`
	LuchadorID       uint `json:"luchadorID"`
}

// ScoreList definition
type ScoreList struct {
	Scores []MatchScore `json:"scores"`
}

// MatchParticipant definition
type MatchParticipant struct {
	LuchadorID uint `json:"luchadorID"`
	MatchID    uint `json:"matchID"`
}

// MatchScore definition
type MatchScore struct {
	ID         uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `json:"-" faker:"-"`
	LuchadorID uint       `json:"luchadorID"`
	MatchID    uint       `json:"matchID"`
	Kills      int        `json:"kills"`
	Deaths     int        `json:"deaths"`
	Score      int        `json:"score"`
}

// MatchMetric definition
type MatchMetric struct {
	ID               uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	MatchID          uint       `json:"matchID"`
	FPS              uint       `json:"fps"`
	Players          uint       `json:"players"`
	GameDefinitionID uint       `json:"gameDefinitionID"`
}

// Classroom definition
type Classroom struct {
	ID         uint       `gorm:"primary_key" json:"id,omitempty" faker:"-"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `json:"-" faker:"-"`
	Name       string     `json:"name"`
	AccessCode string     `json:"accessCode" gorm:"not null;unique_index"`
	OwnerID    uint       `json:"ownerID,omitempty"`
	Students   []Student  `gorm:"many2many:classroom_students" json:"students"`
}

// Student definition
type Student struct {
	ID         uint        `gorm:"primary_key" json:"id,omitempty" faker:"-"`
	CreatedAt  time.Time   `json:"-"`
	UpdatedAt  time.Time   `json:"-"`
	DeletedAt  *time.Time  `json:"-" faker:"-"`
	UserID     uint        `json:"userID,omitempty" gorm:"not null;unique_index"`
	Classrooms []Classroom `gorm:"many2many:classroom_students" json:"classrooms"`
}

// StudentResponse definition
type StudentResponse struct {
	StudentID uint   `json:"studentID"`
	UserID    uint   `json:"userID"`
	Username  string `json:"username"`
}

// AvailableMatch definition
type AvailableMatch struct {
	ID               uint            `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
	DeletedAt        *time.Time      `json:"-" faker:"-"`
	Name             string          `json:"name"`
	GameDefinitionID uint            `json:"gameDefinitionID"`
	ClassroomID      uint            `json:"classroomID"`
	GameDefinition   *GameDefinition `json:"gameDefinition"`
}

//UpdateLuchadorResponse data structure
type UpdateLuchadorResponse struct {
	Errors   []string       `json:"errors"`
	Luchador *GameComponent `json:"luchador"`
}

// PageEventRequest sent from the application
type PageEventRequest struct {
	Page        string `json:"page"`
	Action      string `json:"action"`
	ComponentID string `json:"componentID"`
	AppName     string `json:"appName"`
	AppVersion  string `json:"appVersion"`
	Value1      string `json:"value1"`
	Value2      string `json:"value2"`
	Value3      string `json:"value3"`
}

// PageEvent to be saved after http request inspection
type PageEvent struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-" faker:"-"`
	UserID      uint       `json:"userID"`
	RemoteAddr  string     `json:"remoteAddr"`
	UserAgent   string     `json:"userAgent"`
	Version     string     `json:"version"`
	OSName      string     `json:"OSName"`
	OSVersion   string     `json:"OSVersion"`
	Mobile      bool       `json:"mobile"`
	Tablet      bool       `json:"tablet"`
	Desktop     bool       `json:"desktop"`
	Device      string     `json:"device"`
	Page        string     `json:"page"`
	Action      string     `json:"action"`
	ComponentID string     `json:"componentID"`
	AppName     string     `json:"AppName"`
	AppVersion  string     `json:"AppVersion"`
	Value1      string     `json:"value1"`
	Value2      string     `json:"value2"`
	Value3      string     `json:"value3"`
}
