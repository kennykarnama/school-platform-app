package core

import (
	"context"
	"fmt"
	"time"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	entityStudent "github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	RegisterStudent(ctx context.Context, student *entityStudent.Student, studentClass *entityStudent.StudentClass, createdBy string) error
	SubmitAttendance(ctx context.Context, data []*student.StudentAttendance) error
	ListAttendance(ctx context.Context, academicYearID string, classLabel string, attendanceDate string, userId string) ([]*entityStudent.Aggregate, error)
	DeactivateStudentClass(ctx context.Context, studentClassID string, reason string) error
	StudentClassesAggregate(ctx context.Context, academicYearID, classLabel, userId string) ([]*entityStudent.Aggregate, error)
	StudentAttendances(ctx context.Context, academicYearID string, classLabel string, from *string, to *string, userId string) ([]*entityStudent.StudentAttendance, error)
	StudentClasses(ctx context.Context, academicYearID, classLabel, userId string) ([]*entityStudent.StudentClass, error)
	SaveStudentClasses(ctx context.Context, studentClasses []*entityStudent.StudentClass) error
}

type sqlRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *sqlRepository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) RegisterStudent(ctx context.Context, student *entityStudent.Student, studentClass *entityStudent.StudentClass, createdBy string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(student).Error; err != nil {
			// return any error will rollback
			return err
		}

		// validate

		var results []*entityStudent.Aggregate
		err := tx.Table("student").Select("student.id AS student_id, student.name, student_class.id as student_class_id, student_attendance.id as student_attendance_id, student_attendance.attend as attend, student_attendance.attendance_date::date as attendance_date, student_attendance.attendance_type_id as attendance_type_id").Joins("inner join student_class ON student_class.student_id = student.id").
			Joins("left join student_attendance ON student_attendance.student_class_id = student_class.id").
			Where("academic_year_id = ? AND class_label = ? AND student_class.deleted_at IS NULL", studentClass.AcademicYearID, studentClass.ClassLabel).Order("student_class.created_at ASC").Find(&results).Error
		if err != nil {
			return err
		}

		m := make(map[string]bool)
		for _, result := range results {
			m[result.Name] = true
		}
		if _, ok := m[student.Name]; ok {
			return fmt.Errorf("action=repo.registerStudent err=%v", fmt.Sprintf("duplicate name=%s", student.Name))
		}

		studentClass.StudentID = student.ID
		studentClass.UserID = createdBy

		if err := tx.Create(studentClass).Error; err != nil {
			return err
		}
		// return nil will commit the whole transaction
		return nil
	})
	return err
}

func (r *sqlRepository) SubmitAttendance(ctx context.Context, data []*entityStudent.StudentAttendance) error {
	err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "student_class_id",
			},
			{
				Name: "attendance_date",
			},
		},
		DoUpdates: clause.AssignmentColumns([]string{"attendance_type_id"}),
	}).Create(data).Error
	return err
}

func (r *sqlRepository) ListAttendance(ctx context.Context, academicYearID string, classLabel string, attendanceDate string, userId string) ([]*entityStudent.Aggregate, error) {
	var results []*entityStudent.Aggregate
	err := r.db.Table("student").Select("student.id AS student_id, student.name, student_class.id as student_class_id, "+
		"student_attendance.id as student_attendance_id, student_attendance.attend as attend, "+
		"student_attendance.attendance_date::date as attendance_date, student_attendance.attendance_type_id as attendance_type_id").
		Joins("inner join student_class ON student_class.student_id = student.id AND student_class.user_id = ?", userId).
		Joins("left join student_attendance ON student_attendance.student_class_id = student_class.id AND student_attendance.attendance_date = ?", attendanceDate).
		Where("academic_year_id = ? AND class_label = ? AND student_class.deleted_at IS NULL", academicYearID, classLabel).Order("student_class.created_at ASC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *sqlRepository) DeactivateStudentClass(ctx context.Context, studentClassID string, reason string) error {
	err := r.db.Table("student_class").Where("id = ?", studentClassID).Updates(map[string]interface{}{
		"deleted_at":        time.Now().UTC(),
		"deactivate_reason": reason,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *sqlRepository) StudentAttendances(ctx context.Context, academicYearID string, classLabel string, from *string, to *string, userId string) ([]*entityStudent.StudentAttendance, error) {
	var results []*entityStudent.StudentAttendance
	q := r.db.Table("student_attendance").Select("student_attendance.*").Joins("inner join student_class on student_class.id = student_attendance.student_class_id and student_class.deleted_at is null")
	if from != nil {
		q = q.Where("student_attendance.attendance_date::date >= ?", *from)
	}
	if to != nil {
		q = q.Where("student_attendance.attendance_date::date <= ?", *to)
	}
	q = q.Where("student_class.academic_year_id = ?", academicYearID)
	q = q.Where("student_class.class_label = ?", classLabel)
	q = q.Where("student_class.user_id = ?", userId)
	err := q.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *sqlRepository) StudentClassesAggregate(ctx context.Context, academicYearID, classLabel, userId string) ([]*entityStudent.Aggregate, error) {
	var results []*entityStudent.Aggregate
	err := r.db.Table("student_class").Select("student.id AS student_id, student.name, student_class.id as student_class_id").
		Joins("inner join student on student.id = student_class.student_id").
		Where("academic_year_id = ? AND class_label = ? AND student_class.deleted_at IS NULL AND student_class.user_id = ?", academicYearID, classLabel, userId).
		Order("student_class.created_at ASC").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *sqlRepository) StudentClasses(ctx context.Context, academicYearID, classLabel, userId string) ([]*entityStudent.StudentClass, error) {
	var results []*entityStudent.StudentClass
	err := r.db.Where("academic_year_id = ? AND class_label = ? AND student_class.deleted_at IS NULL AND student_class.user_id = ?", academicYearID, classLabel, userId).
		Order("student_class.created_at ASC").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *sqlRepository) SaveStudentClasses(ctx context.Context, studentClasses []*entityStudent.StudentClass) error {
	err := r.db.Save(studentClasses).Error
	if err != nil {
		return err
	}
	return nil
}
