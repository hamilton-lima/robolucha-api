package main

import "time"

// User definition
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Username  string     `json:"username"`
}

// Session definition
type Session struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	UUID      string     `json:"UUID"`
	UserID    uint       `json:"userID"`
}

// UserSetting definition
type UserSetting struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `json:"-"`
	UserID     uint       `json:"userID"`
	LastOption string     `json:"lastOption"`
}

// Match definition
type Match struct {
	ID               uint            `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
	DeletedAt        *time.Time      `json:"-"`
	TimeStart        time.Time       `json:"timeStart"`
	TimeEnd          time.Time       `json:"timeEnd"`
	LastTimeAlive    time.Time       `json:"lastTimeAlive"`
	GameDefinitionID uint            `json:"gameDefinitionID"`
	Participants     []GameComponent `gorm:"many2many:match_participants" json:"participants"`
}

type GameDefinition struct {
	ID                            uint             `gorm:"primary_key" json:"id"`
	CreatedAt                     time.Time        `json:"-"`
	UpdatedAt                     time.Time        `json:"-"`
	DeletedAt                     *time.Time       `json:"-"`
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
	GameComponents                []GameComponent  `json:"gameComponents"`
	SceneComponents               []SceneComponent `json:"sceneComponents"`
	Codes                         []Code           `gorm:"many2many:gamedefinition_codes" json:"codes"`
	LuchadorSuggestedCodes        []Code           `gorm:"many2many:gamedefinition_suggestedcodes" json:"suggestedCodes"`
}

type SceneComponent struct {
	ID               uint       `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-"`
	GameDefinitionID uint       `json:"gameDefinition,omitempty"`
	X                uint       `json:"x"`
	Y                uint       `json:"y"`
	Width            uint       `json:"width"`
	Height           uint       `json:"height"`
	Rotation         uint       `json:"rotation"`
	Respawn          bool       `json:"respawn"`
	Colider          bool       `json:"colider"`
	ShowInRadar      bool       `json:"showInRadar"`
	BlockMovement    bool       `json:"blockMovement"`
	Type             string     `json:"name"`
	Codes            []Code     `gorm:"many2many:scenecomponent_codes" json:"codes"`
}

type GameComponent struct {
	ID               uint       `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-"`
	GameDefinitionID uint       `json:"gameDefinition,omitempty"`
	Name             string     `gorm:"not null;unique_index" json:"name"`
	UserID           uint       `json:"userID,omitempty"`
	Codes            []Code     `gorm:"many2many:gamecomponent_codes" json:"codes"`
	Configs          []Config   `gorm:"many2many:gamecomponent_configs" json:"configs"`
}

// Code definition
type Code struct {
	ID               uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-"`
	Event            string     `json:"event"`
	Script           string     `json:"script"`
	Exception        string     `json:"exception"`
	GameDefinitionID uint       `json:"gameDefinition"`
}

// Config definition
type Config struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Key       string     `json:"key"`
	Value     string     `json:"value"`
}

// JoinMatch definition
type JoinMatch struct {
	MatchID    uint `json:"matchID"`
	LuchadorID uint `json:"luchadorID"`
}

// JoinMatch definition
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
	DeletedAt  *time.Time `json:"-"`
	LuchadorID uint       `json:"luchadorID"`
	MatchID    uint       `json:"matchID"`
	Kills      int        `json:"kills"`
	Deaths     int        `json:"deaths"`
	Score      int        `json:"score"`
}
