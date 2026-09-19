package class

import (
	"context"
	classEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/class"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, principal user.Principal) ([]*classEntity.Class, error)
}

type sqlRepository struct{ db *gorm.DB }

func NewSQLRepository(db *gorm.DB) *sqlRepository { return &sqlRepository{db: db} }

func (r *sqlRepository) List(ctx context.Context, principal user.Principal) ([]*classEntity.Class, error) {
	var results []*classEntity.Class
	query := r.db.WithContext(ctx).Where("school_id = ? AND active = true", principal.SchoolID)
	if principal.Role == user.RoleTeacher {
		query = query.Where("EXISTS (SELECT 1 FROM teacher_class_access access WHERE access.school_id = school_class.school_id AND access.class_label = school_class.label AND access.teacher_id = ?)", principal.UserID)
	}
	return results, query.Order("label ASC").Find(&results).Error
}
