package academicyear

import (
	"context"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/academicyear"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, principal user.Principal) ([]*academicyear.AcademicYear, error)
}

type sqlRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *sqlRepository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) List(ctx context.Context, principal user.Principal) ([]*academicyear.AcademicYear, error) {
	var results []*academicyear.AcademicYear
	query := r.db.WithContext(ctx).Where("school_id = ?", principal.SchoolID)
	if principal.Role == user.RoleTeacher {
		query = query.Where("EXISTS (SELECT 1 FROM teacher_class_access access WHERE access.academic_year_id = academic_year.id AND access.teacher_id = ? AND access.school_id = ?)", principal.UserID, principal.SchoolID)
	}
	err := query.Order("created_at asc").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
