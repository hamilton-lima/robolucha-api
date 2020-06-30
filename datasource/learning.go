package datasource

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/model"
)

// FindAllActivities definition
func (ds *DataSource) FindAllActivities() *[]model.Activity {
	var result []model.Activity

	ds.DB.
		Preload("Skill").
		Find(&result)

	log.WithFields(log.Fields{
		"activities": result,
	}).Debug("FindAllActivities")

	return &result
}

// FindLearningObjectiveByName definition
func (ds *DataSource) FindLearningObjectiveByName(name string) *model.LearningObjective {
	var result model.LearningObjective

	if ds.DB.Where(&model.LearningObjective{Name: name}).Find(&result).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"learning objective": result,
	}).Debug("FindLearningObjectiveByName")

	return &result
}

// AddLearningObjective definition
func (ds *DataSource) AddLearningObjective(c *model.LearningObjective) *model.LearningObjective {

	objective := model.LearningObjective{
		Name:   c.Name,
		Skills: c.Skills,
	}

	log.WithFields(log.Fields{
		"objective": objective,
	}).Debug("AddLearningObjective")

	ds.DB.Create(&objective)

	log.WithFields(log.Fields{
		"objective": objective,
	}).Debug("after AddLearningObjective")

	return &objective
}



// FindAllAssignment definition
func (ds *DataSource) FindAllAssignments() *[]model.Assignment {
	ds.DB.LogMode(true)
	var result []model.Assignment

	ds.DB.
		Preload("Students").
		Preload("Activities").
		Find(&result)

	log.WithFields(log.Fields{
		"assignments": result,
	}).Debug("FindAllAssignments")

	return &result
}

// FindAllAssignment definition
func (ds *DataSource) FindAssignmentById(id uint) *model.Assignment {
	ds.DB.LogMode(true)
	var result model.Assignment
	ds.DB.
		Preload("Students").
		Preload("Activities").
		Where(&model.Assignment{ID: id}).
		Find(&result)

	log.WithFields(log.Fields{
		"assignments": result,
	}).Debug("FindAssignmentById")

	return &result
}

func (ds *DataSource) AddAssignment(assignment *model.Assignment) *model.Assignment {
	newAssignment := model.Assignment{
		TimeStart: assignment.TimeStart,
		TimeEnd:   assignment.TimeEnd,
	}

	ds.DB.Create(&newAssignment)

	return &newAssignment
}

func (ds *DataSource) DeleteAssignment(id uint) {
	ds.DB.LogMode(true)
	ds.DB.
		Where(&model.Assignment{ID: id}).
		Delete(model.Assignment{})

	log.WithFields(log.Fields{
		"assignment.id": id,
	}).Debug("DeleteAssignment")

}

func (ds *DataSource) UpdateAssignmentStudents(id uint, studentIds []uint) *model.Assignment {
	var assignment model.Assignment
	ds.DB.LogMode(true)
	ds.DB.
		Preload("Students").
		Where(&model.Assignment{ID: id}).
		Find(&assignment)

	students := make([]model.Student, 0)
	for _, studentId := range studentIds {
		var student model.Student
		ds.DB.
			Where(&model.Student{ID: studentId}).
			Find(&student)

		students = append(students, student)
	}
	assignment.Students = students
	ds.DB.Save(&assignment)
	return &assignment
}

func (ds *DataSource) UpdateAssignmentActivities(id uint, activityIds []uint) *model.Assignment {
	var assignment model.Assignment
	ds.DB.LogMode(true)
	ds.DB.
		Preload("Activities").
		Where(&model.Assignment{ID: id}).
		Find(&assignment)

	activities := make([]model.Activity, 0)
	for _, activityId := range activityIds {
		var activity model.Activity
		ds.DB.
			Where(&model.Activity{ID: activityId}).
			Find(&activity)

		activities = append(activities, activity)
	}
	assignment.Activities = activities
	ds.DB.Save(&assignment)
	return &assignment
}

