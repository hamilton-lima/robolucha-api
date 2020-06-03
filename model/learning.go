package model

import "time"

// LearningObjective definition
type LearningObjective struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" faker:"-"`
	Name      string     `json:"name" gorm:"not null;unique_index"`
	Skills    []Skill    `gorm:"many2many:learningobjective_skills" json:"skills"`
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

// LevelGroup definition
type LevelGroup struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-" faker:"-"`
	Name        string     `json:"name" gorm:"not null;unique_index"`
	Description string     `gorm:"size:125000" json:"description"`
	MinLevel    uint       `json:"minLevel"`
}

// TODO add this as second phase
// // GradingSystem definition
// type GradingSystem struct {
// 	ID        uint       `gorm:"primary_key" json:"id"`
// 	CreatedAt time.Time  `json:"-"`
// 	UpdatedAt time.Time  `json:"-"`
// 	DeletedAt *time.Time `json:"-" faker:"-"`
// 	Name      string     `json:"name"`
// 	Grades    []Grade    `json:"grades"`
// }

// TODO: Add Grade
// GradingSystemID uint       `json:"gradingSystemID"`

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

// TODO: Add to Activity
// GradingSystemID  uint            `json:"gradingSystemID"`
// GradingSystem    GradingSystem   `json:"gradingSystem"`

// Activity definition
type Activity struct {
	ID               uint            `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
	DeletedAt        *time.Time      `json:"-" faker:"-"`
	Name             string          `json:"name"`
	Description      string          `gorm:"size:125000" json:"description"`
	Skills           []Skill         `gorm:"many2many:activitiy_skills" json:"skills"`
	GameDefinitionID uint            `json:"gameDefinitionID"`
	GameDefinition   *GameDefinition `json:"gameDefinition"`
	SourceURL        string          `json:"sourceURL"`
	SourceName       string          `json:"sourceName"`
}

// Assignment definition
type Assignment struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-" faker:"-"`
	ClassroomID uint       `json:"classroomID"`
	Classroom   Classroom  `json:"classroom"`
	ActivityID  uint       `json:"activityID"`
	Activity    Activity   `json:"activity"`
	TimeStart   time.Time  `json:"timeStart"`
	TimeEnd     time.Time  `json:"timeEnd"`
}

// AssignmentEvaluation definition
type AssignmentEvaluation struct {
	ID               uint              `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time         `json:"-"`
	UpdatedAt        time.Time         `json:"-"`
	DeletedAt        *time.Time        `json:"-" faker:"-"`
	AssignmentID     uint              `json:"assignmentID"`
	Assignment       Assignment        `json:"assignment"`
	StudentID        uint              `json:"studentID"`
	Student          Student           `json:"student"`
	AssignmentGrades []AssignmentGrade `json:"assignmentGrades"`
}

// AssignmentGrade definition
type AssignmentGrade struct {
	ID                     uint       `gorm:"primary_key" json:"id"`
	CreatedAt              time.Time  `json:"-"`
	UpdatedAt              time.Time  `json:"-"`
	DeletedAt              *time.Time `json:"-" faker:"-"`
	Grade                  float32    `json:"grade"`
	SkillID                uint       `json:"skillID"`
	Skill                  Skill      `json:"skill"`
	AssignmentEvaluationID uint       `json:"assignmentEvaluationID"`
}
