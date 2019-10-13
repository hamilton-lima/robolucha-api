package model

import "time"

// LearningObjective definition
type LearningObjective struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	Name      string     `json:"name"`
}

// Skill definition
type Skill struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-" faker:"-"`
	Name        string     `json:"name" gorm:"not null;unique_index"`
	Description string     `gorm:"size:125000" json:"description"`
}

// GradingSystem definition
type GradingSystem struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	Name      string     `json:"name"`
	Grades    []Grade    `json:"grades"`
}

// Grade definition
type Grade struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	Name      string     `json:"name"`
	Lowest    float32    `json:"lowest"`
	Highest   float32    `json:"highest"`
	Color     string     `json:"color"`
}

// Check canada learning code data structure for
// activities adn mine craft education
// check license of the activities
// add activities without gamedefinition???

// Activity definition
type Activity struct {
	ID               uint       `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-" faker:"-"`
	Name             string     `json:"name"`
	Description      string     `gorm:"size:125000" json:"description"`
	Skills           []Skill    `gorm:"many2many:activitiy_skills" json:"skills"`
	GameDefinitionID uint       `json:"gameDefinitionID"`
	SourceURL        string     `json:"sourceURL"`
	SourceName       string     `json:"sourceName"`
}
