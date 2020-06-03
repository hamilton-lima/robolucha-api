package datasource

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/model"
)

// TODO: remove this
var accessCodeCounter int64

func (ds *DataSource) AddClassroom(c *model.Classroom) *model.Classroom {

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

	ds.DB.Create(&classroom)
	classroom.Students = make([]model.Student, 0)

	log.WithFields(log.Fields{
		"classroom": classroom,
	}).Debug("after addClassroom")

	// create avaialable match for all existing gamedefinitions
	all := ds.FindAllGameDefinition()
	for _, gd := range *all {

		am := model.AvailableMatch{GameDefinitionID: gd.ID, Name: gd.Name, ClassroomID: classroom.ID}
		ds.DB.Where(&am).FirstOrCreate(&am)

		log.WithFields(log.Fields{
			"gameDefinitionID": gd.ID,
			"AvailableMatch":   am,
			"classroom":        classroom,
		}).Debug("AddClassroom")
	}

	return &classroom
}

func (ds *DataSource) FindAllClassroom(user *model.User) *[]model.Classroom {
	var result []model.Classroom

	ds.DB.
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

func (ds *DataSource) FindClassroomByID(id uint) *model.Classroom {
	var result model.Classroom

	ds.DB.
		Preload("Students").
		Where(&model.Classroom{ID: id}).
		Find(&result)

	log.WithFields(log.Fields{
		"classroom": result,
	}).Debug("FindClassroomByID")

	return &result
}

// FindAllClassroomByStudent definition
func (ds *DataSource) FindAllClassroomByStudent(studentUserID uint) []model.Classroom {

	student := model.Student{UserID: studentUserID}

	if ds.DB.
		Preload("Classrooms").
		Where(&student).
		First(&student).
		RecordNotFound() {

		log.WithFields(log.Fields{
			"userID": studentUserID,
		}).Error("student not found")

		return make([]model.Classroom, 0)
	}

	log.WithFields(log.Fields{
		"classrooms": student.Classrooms,
	}).Debug("FindAllClassroomByStudent")

	return student.Classrooms
}

// JoinClassroom definition
func (ds *DataSource) JoinClassroom(user *model.User, accessCode string) *model.Classroom {
	var result model.Classroom
	student := model.Student{UserID: user.ID}

	if ds.DB.Preload("Students").
		Where(&model.Classroom{AccessCode: accessCode}).
		First(&result).
		RecordNotFound() {

		log.WithFields(log.Fields{
			"accessCode": accessCode,
		}).Info("classroom not found")

		return nil
	}

	if ds.DB.
		Where(&student).
		First(&student).
		RecordNotFound() {

		log.WithFields(log.Fields{
			"userID": user.ID,
		}).Info("student not found will create")

		ds.DB.Create(&student)
	}

	result.Students = append(result.Students, student)
	ds.DB.Save(&result)

	log.WithFields(log.Fields{
		"classroom": result,
		"student":   student,
	}).Info("joinClassroom")

	return &result
}

func (ds *DataSource) FindPublicAvailableMatch() *[]model.AvailableMatch {
	return ds.FindAvailableMatchByClassroomID(0)
}

func (ds *DataSource) FindAvailableMatchByClassroomID(id uint) *[]model.AvailableMatch {

	var result []model.AvailableMatch

	ds.DB.
		Where("classroom_id = ?", id).
		Find(&result)

	log.WithFields(log.Fields{
		"availableMatch": result,
	}).Debug("findAvailableMatchByClassroomID")

	// load game definition details
	for n, availableMatch := range result {
		result[n].GameDefinition = ds.FindGameDefinition(availableMatch.GameDefinitionID)
	}

	return &result
}

// BuildStudentResponse definition
func (ds *DataSource) BuildStudentResponse(students []model.Student) []model.StudentResponse {
	result := make([]model.StudentResponse, len(students))

	for n, student := range students {
		user := ds.FindUserByID(student.UserID)
		result[n] = model.StudentResponse{
			StudentID: student.ID,
			UserID:    student.UserID,
			Username:  user.Username,
		}
	}

	return result
}

// FindGradeByName definition
func (ds *DataSource) FindGradeByName(name string) *model.Grade {
	var result model.Grade

	if ds.DB.Where(&model.Grade{Name: name}).Find(&result).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"grade": result,
	}).Debug("FindGradeByName")

	return &result
}

// FindLevelGroupByName definition
func (ds *DataSource) FindLevelGroupByName(name string) *model.LevelGroup {
	var result model.LevelGroup

	if ds.DB.Where(&model.LevelGroup{Name: name}).Find(&result).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"level group": result,
	}).Debug("FindLevelGroupByName")

	return &result
}

// AddGrade defines add grade
func (ds *DataSource) AddGrade(c *model.Grade) *model.Grade {

	grade := model.Grade{
		Name:    c.Name,
		Lowest:  c.Lowest,
		Highest: c.Highest,
		Color:   c.Color,
	}

	log.WithFields(log.Fields{
		"grade": grade,
	}).Debug("AddGrade")

	ds.DB.Create(&grade)

	log.WithFields(log.Fields{
		"classroom": grade,
	}).Debug("after AddGrade")

	return &grade
}

// AddLevelGroup definition
func (ds *DataSource) AddLevelGroup(c *model.LevelGroup) *model.LevelGroup {

	levelGroup := model.LevelGroup{
		Name:        c.Name,
		Description: c.Description,
		MinLevel:    c.MinLevel,
	}

	log.WithFields(log.Fields{
		"levelGroup": levelGroup,
	}).Debug("AddLevelGroup")

	ds.DB.Create(&levelGroup)

	log.WithFields(log.Fields{
		"levelGroup": levelGroup,
	}).Debug("after AddLevelGroup")

	return &levelGroup
}

// FindLevelGroup definition
func (ds *DataSource) FindLevelGroup() *[]model.LevelGroup {
	var result []model.LevelGroup

	ds.DB.Find(&result)

	log.WithFields(log.Fields{
		"levelgroups": result,
	}).Debug("FindLevelGroup")

	return &result
}
