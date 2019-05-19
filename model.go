package main

import "time"

// User definition
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Username  string     `json:"username"`
}

// Session definition
type Session struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	UUID      string     `json:"UUID"`
	UserID    uint       `json:"userID"`
}

// UserSetting definition
type UserSetting struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
	UserID     uint       `json:"userID"`
	LastOption string     `json:"lastOption"`
}

// Match definition
type Match struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt"`
	TimeStart       time.Time  `json:"timeStart"`
	TimeEnd         time.Time  `json:"timeEnd"`
	LastTimeAlive   time.Time  `json:"lastTimeAlive"`
	Duration        uint64     `json:"duration"`
	Participants    []Luchador `gorm:"many2many:match_participants" json:"participants"`
	MinParticipants uint       `json:"minParticipants"`
	MaxParticipants uint       `json:"maxParticipants"`
	ArenaWidth      uint       `json:"arenaWidth"`
	ArenaHeight     uint       `json:"arenaHeight"`
	BulletSize      uint       `json:"bulletSize"`
	LuchadorSize    uint       `json:"luchadorSize"`
	Fps             uint       `json:"fps"`
	BuletSpeed      uint       `json:"buletSpeed"`
}

type GameDefinition struct {
	ID                     uint             `gorm:"primary_key" json:"id"`
	CreatedAt              time.Time        `json:"createdAt"`
	UpdatedAt              time.Time        `json:"updatedAt"`
	DeletedAt              *time.Time       `json:"deletedAt"`
	Duration               uint64           `json:"duration"`
	MinParticipants        uint             `json:"minParticipants"`
	MaxParticipants        uint             `json:"maxParticipants"`
	ArenaWidth             uint             `json:"arenaWidth"`
	ArenaHeight            uint             `json:"arenaHeight"`
	BulletSize             uint             `json:"bulletSize"`
	LuchadorSize           uint             `json:"luchadorSize"`
	Fps                    uint             `json:"fps"`
	BuletSpeed             uint             `json:"buletSpeed"`
	Name                   string           `gorm:"not null;unique_index" json:"name"`
	Description            string           `json:"description"`
	Type                   string           `json:"type"`
	SortOrder              uint             `json:"sortOrder"`
	Participants           []Luchador       `gorm:"many2many:match_participants" json:"participants"`
	SceneComponents        []SceneComponent `gorm:"many2many:scene_components" json:"scene_components"`
	Codes                  []Code           `json:"codes"`
	LuchadorSuggestedCodes []Code           `json:"suggestedCodes"`
}

// Luchador definition
type SceneComponent struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
	X             uint       `json:"x"`
	Y             uint       `json:"y"`
	Width         uint       `json:"width"`
	Height        uint       `json:"height"`
	Rotation      uint       `json:"rotation"`
	Respawn       bool       `json:"respawn"`
	Colider       bool       `json:"colider"`
	ShowInRadar   bool       `json:"showInRadar"`
	BlockMovement bool       `json:"blockMovement"`
	Type          string     `json:"name"`
	Codes         []Code     `json:"codes"`
}

// Luchador definition
type Luchador struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	UserID    uint       `json:"userID"`
	Name      string     `gorm:"not null;unique_index" json:"name"`
	Codes     []Code     `json:"codes"`
	Configs   []Config   `json:"configs"`
}

// Code definition
type Code struct {
	ID         uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt  time.Time  `json:"createdAt,omitempty"`
	UpdatedAt  time.Time  `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	LuchadorID uint       `json:"luchadorID,omitempty"`
	Event      string     `json:"event"`
	Script     string     `json:"script"`
	Exception  string     `json:"exception"`
}

// Config definition
type Config struct {
	ID         uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt  time.Time  `json:"createdAt,omitempty"`
	UpdatedAt  time.Time  `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	LuchadorID uint       `json:"luchadorID,omitempty"`
	Key        string     `json:"key"`
	Value      string     `json:"value"`
}

// JoinMatch definition
type JoinMatch struct {
	MatchID    uint `json:"matchID"`
	LuchadorID uint `json:"luchadorID"`
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
	CreatedAt  time.Time  `json:"createdAt,omitempty"`
	UpdatedAt  time.Time  `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	LuchadorID uint       `json:"luchadorID"`
	MatchID    uint       `json:"matchID"`
	Kills      int        `json:"kills"`
	Deaths     int        `json:"deaths"`
	Score      int        `json:"score"`
}
