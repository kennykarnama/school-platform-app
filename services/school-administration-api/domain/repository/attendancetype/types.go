package attendancetype

import (
	"context"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/attendancetype"
	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context) ([]*attendancetype.AttendanceType, error)
}

type sqlRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *sqlRepository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) List(ctx context.Context) ([]*attendancetype.AttendanceType, error) {
	var results []*attendancetype.AttendanceType
	err := r.db.Order("created_at ASC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
