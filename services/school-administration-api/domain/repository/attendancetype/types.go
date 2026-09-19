package attendancetype

import (
	"context"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/attendancetype"
	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, schoolID string) ([]*attendancetype.AttendanceType, error)
}

type sqlRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *sqlRepository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) List(ctx context.Context, schoolID string) ([]*attendancetype.AttendanceType, error) {
	var results []*attendancetype.AttendanceType
	err := r.db.WithContext(ctx).Where("school_id = ?", schoolID).Order("created_at ASC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
