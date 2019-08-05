package model

import "time"

// User definition
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	Username  string     `json:"username"`
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
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `json:"-" faker:"-"`
	UserID     uint       `json:"userID"`
	LastOption string     `json:"lastOption"`
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

// Match definition
type Match struct {
	ID               uint            `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
	DeletedAt        *time.Time      `json:"-" faker:"-"`
	TimeStart        time.Time       `json:"timeStart"`
	TimeEnd          time.Time       `json:"timeEnd"`
	LastTimeAlive    time.Time       `json:"lastTimeAlive"`
	GameDefinitionID uint            `json:"gameDefinitionID"`
	Participants     []GameComponent `gorm:"many2many:match_participants" json:"participants"`
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
	GameComponents                []GameComponent  `json:"gameComponents"`
	SceneComponents               []SceneComponent `json:"sceneComponents"`
	Codes                         []Code           `gorm:"many2many:gamedefinition_codes" json:"codes"`
	LuchadorSuggestedCodes        []Code           `gorm:"many2many:gamedefinition_suggestedcodes" json:"suggestedCodes"`
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
}

// Code definition
type Code struct {
	ID               uint       `gorm:"primary_key" json:"id,omitempty" faker:"-"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	Event            string     `json:"event"`
	Script           string     `json:"script"`
	Exception        string     `json:"exception"`
	GameDefinitionID uint       `json:"gameDefinition" faker:"-"`
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
	ID        uint       `gorm:"primary_key" json:"id,omitempty" faker:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	UserID    uint       `json:"userID,omitempty" gorm:"not null;unique_index"`
}

// AvailableMatch definition
type AvailableMatch struct {
	ID               uint       `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	Name             string     `json:"name"`
	GameDefinitionID uint       `json:"gameDefinitionID"`
	Classroom        *Classroom `gorm:"many2many:available_match_classrooms" json:"classroom"`
}

//UpdateLuchadorResponse data structure
type UpdateLuchadorResponse struct {
	Errors   []string       `json:"errors"`
	Luchador *GameComponent `json:"luchador"`
}
