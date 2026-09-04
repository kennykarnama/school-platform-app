package core

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	entityStudent "github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	RegisterStudent(ctx context.Context, student *entityStudent.Student, studentClass *entityStudent.StudentClass, principal user.Principal) error
	SubmitAttendance(ctx context.Context, data []*student.StudentAttendance, principal user.Principal) error
	ListAttendance(ctx context.Context, academicYearID string, classLabel string, attendanceDate string, principal user.Principal) ([]*entityStudent.Aggregate, error)
	DeactivateStudentClass(ctx context.Context, studentClassID string, reason string, principal user.Principal) error
	StudentClassesAggregate(ctx context.Context, academicYearID, classLabel string, principal user.Principal) ([]*entityStudent.Aggregate, error)
	StudentAttendances(ctx context.Context, academicYearID string, classLabel string, from *string, to *string, principal user.Principal) ([]*entityStudent.StudentAttendance, error)
	StudentClasses(ctx context.Context, academicYearID, classLabel string, principal user.Principal) ([]*entityStudent.StudentClass, error)
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

func (r *sqlRepository) RegisterStudent(ctx context.Context, student *entityStudent.Student, studentClass *entityStudent.StudentClass, principal user.Principal) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if !hasSchoolClassAndYear(tx, principal.SchoolID, studentClass.AcademicYearID, studentClass.ClassLabel) {
			return user.ErrForbidden
		}
		if principal.Role == user.RoleTeacher && !hasClassAccess(tx, principal, studentClass.AcademicYearID, studentClass.ClassLabel) {
			return user.ErrForbidden
		}
		student.SchoolID = principal.SchoolID
		if err := tx.Create(student).Error; err != nil {
			return err
		}

		studentClass.StudentID = student.ID
		studentClass.SchoolID = principal.SchoolID
		studentClass.UserID = principal.UserID

		if err := tx.Create(studentClass).Error; err != nil {
			return err
		}
		return nil
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return entityStudent.ErrAlternativeIDAlreadyExists
	}
	return err
}

func (r *sqlRepository) SubmitAttendance(ctx context.Context, data []*entityStudent.StudentAttendance, principal user.Principal) error {
	if len(data) == 0 {
		return nil
	}
	ids := make([]string, 0, len(data))
	typeIDs := make([]string, 0, len(data))
	uniqueStudentClasses := map[string]bool{}
	for _, item := range data {
		ids = append(ids, item.StudentClassID)
		typeIDs = append(typeIDs, item.AttendanceTypeID)
		uniqueStudentClasses[item.StudentClassID] = true
		item.SchoolID = principal.SchoolID
	}
	var allowed int64
	query := r.db.WithContext(ctx).Table("student_class").Where("id IN ? AND school_id = ? AND deleted_at IS NULL", ids, principal.SchoolID)
	if principal.Role == user.RoleTeacher {
		query = withTeacherAccess(query, principal)
	}
	if err := query.Count(&allowed).Error; err != nil {
		return err
	}
	if allowed != int64(len(uniqueStudentClasses)) {
		return user.ErrForbidden
	}
	var allowedTypes int64
	if err := r.db.WithContext(ctx).Table("attendance_type").Where("id IN ? AND school_id = ? AND deleted_at IS NULL", typeIDs, principal.SchoolID).Count(&allowedTypes).Error; err != nil {
		return err
	}
	uniqueTypes := map[string]bool{}
	for _, id := range typeIDs {
		uniqueTypes[id] = true
	}
	if allowedTypes != int64(len(uniqueTypes)) {
		return user.ErrForbidden
	}
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
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

func (r *sqlRepository) ListAttendance(ctx context.Context, academicYearID string, classLabel string, attendanceDate string, principal user.Principal) ([]*entityStudent.Aggregate, error) {
	if !hasSchoolClassAndYear(r.db.WithContext(ctx), principal.SchoolID, academicYearID, classLabel) || (principal.Role == user.RoleTeacher && !hasClassAccess(r.db.WithContext(ctx), principal, academicYearID, classLabel)) {
		return nil, user.ErrForbidden
	}
	var results []*entityStudent.Aggregate
	query := r.db.WithContext(ctx).Table("student").Select("student.id AS student_id, student.name, student_class.id as student_class_id, "+
		"student_attendance.id as student_attendance_id, student_attendance.attend as attend, "+
		"student_attendance.attendance_date::date as attendance_date, student_attendance.attendance_type_id as attendance_type_id").
		Joins("inner join student_class ON student_class.student_id = student.id AND student_class.school_id = ?", principal.SchoolID).
		Joins("left join student_attendance ON student_attendance.student_class_id = student_class.id AND student_attendance.attendance_date = ?", attendanceDate).
		Where("student.school_id = ? AND academic_year_id = ? AND class_label = ? AND student_class.deleted_at IS NULL", principal.SchoolID, academicYearID, classLabel)
	if principal.Role == user.RoleTeacher {
		query = withTeacherAccess(query, principal)
	}
	err := query.Order("student_class.created_at ASC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *sqlRepository) DeactivateStudentClass(ctx context.Context, studentClassID string, reason string, principal user.Principal) error {
	query := r.db.WithContext(ctx).Table("student_class").Where("id = ? AND school_id = ?", studentClassID, principal.SchoolID)
	if principal.Role == user.RoleTeacher {
		query = withTeacherAccess(query, principal)
	}
	result := query.Updates(map[string]interface{}{
		"deleted_at":        time.Now().UTC(),
		"deactivate_reason": reason,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return user.ErrForbidden
	}
	return nil
}

func (r *sqlRepository) StudentAttendances(ctx context.Context, academicYearID string, classLabel string, from *string, to *string, principal user.Principal) ([]*entityStudent.StudentAttendance, error) {
	if !hasSchoolClassAndYear(r.db.WithContext(ctx), principal.SchoolID, academicYearID, classLabel) || (principal.Role == user.RoleTeacher && !hasClassAccess(r.db.WithContext(ctx), principal, academicYearID, classLabel)) {
		return nil, user.ErrForbidden
	}
	var results []*entityStudent.StudentAttendance
	q := r.db.WithContext(ctx).Table("student_attendance").Select("student_attendance.*").Joins("inner join student_class on student_class.id = student_attendance.student_class_id and student_class.deleted_at is null").Where("student_attendance.school_id = ? AND student_class.school_id = ?", principal.SchoolID, principal.SchoolID)
	if from != nil {
		q = q.Where("student_attendance.attendance_date::date >= ?", *from)
	}
	if to != nil {
		q = q.Where("student_attendance.attendance_date::date <= ?", *to)
	}
	q = q.Where("student_class.academic_year_id = ?", academicYearID)
	q = q.Where("student_class.class_label = ?", classLabel)
	if principal.Role == user.RoleTeacher {
		q = withTeacherAccess(q, principal)
	}
	err := q.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *sqlRepository) StudentClassesAggregate(ctx context.Context, academicYearID, classLabel string, principal user.Principal) ([]*entityStudent.Aggregate, error) {
	if !hasSchoolClassAndYear(r.db.WithContext(ctx), principal.SchoolID, academicYearID, classLabel) || (principal.Role == user.RoleTeacher && !hasClassAccess(r.db.WithContext(ctx), principal, academicYearID, classLabel)) {
		return nil, user.ErrForbidden
	}
	var results []*entityStudent.Aggregate
	query := r.db.WithContext(ctx).Table("student_class").Select("student.id AS student_id, student.name, student_class.id as student_class_id").
		Joins("inner join student on student.id = student_class.student_id").
		Where("student.school_id = ? AND student_class.school_id = ? AND academic_year_id = ? AND class_label = ? AND student_class.deleted_at IS NULL", principal.SchoolID, principal.SchoolID, academicYearID, classLabel)
	if principal.Role == user.RoleTeacher {
		query = withTeacherAccess(query, principal)
	}
	err := query.Order("student_class.created_at ASC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *sqlRepository) StudentClasses(ctx context.Context, academicYearID, classLabel string, principal user.Principal) ([]*entityStudent.StudentClass, error) {
	if !hasSchoolClassAndYear(r.db.WithContext(ctx), principal.SchoolID, academicYearID, classLabel) || (principal.Role == user.RoleTeacher && !hasClassAccess(r.db.WithContext(ctx), principal, academicYearID, classLabel)) {
		return nil, user.ErrForbidden
	}
	var results []*entityStudent.StudentClass
	query := r.db.WithContext(ctx).Where("school_id = ? AND academic_year_id = ? AND class_label = ? AND student_class.deleted_at IS NULL", principal.SchoolID, academicYearID, classLabel)
	if principal.Role == user.RoleTeacher {
		query = withTeacherAccess(query, principal)
	}
	err := query.Order("student_class.created_at ASC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func withTeacherAccess(query *gorm.DB, principal user.Principal) *gorm.DB {
	return query.Where("EXISTS (SELECT 1 FROM teacher_class_access access WHERE access.school_id = student_class.school_id AND access.academic_year_id = student_class.academic_year_id AND access.class_label = student_class.class_label AND access.teacher_id = ?)", principal.UserID)
}

func hasClassAccess(db *gorm.DB, principal user.Principal, academicYearID, classLabel string) bool {
	var count int64
	err := db.Table("teacher_class_access").Where("school_id = ? AND teacher_id = ? AND academic_year_id = ? AND class_label = ?", principal.SchoolID, principal.UserID, academicYearID, classLabel).Count(&count).Error
	return err == nil && count > 0
}

func hasSchoolClassAndYear(db *gorm.DB, schoolID, academicYearID, classLabel string) bool {
	var yearCount, classCount int64
	if err := db.Table("academic_year").Where("id = ? AND school_id = ? AND deleted_at IS NULL", academicYearID, schoolID).Count(&yearCount).Error; err != nil {
		return false
	}
	if err := db.Table("school_class").Where("school_id = ? AND label = ? AND active = true", schoolID, classLabel).Count(&classCount).Error; err != nil {
		return false
	}
	return yearCount == 1 && classCount == 1
}

func (r *sqlRepository) SaveStudentClasses(ctx context.Context, studentClasses []*entityStudent.StudentClass) error {
	err := r.db.Save(studentClasses).Error
	if err != nil {
		return err
	}
	return nil
}
