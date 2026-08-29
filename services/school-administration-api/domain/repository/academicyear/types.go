package academicyear

import (
	"context"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/academicyear"
	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context) ([]*academicyear.AcademicYear, error)
}

type sqlRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *sqlRepository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) List(ctx context.Context) ([]*academicyear.AcademicYear, error) {
	var results []*academicyear.AcademicYear
	err := r.db.Order("created_at asc").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
