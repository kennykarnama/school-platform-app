package attendancetype

import (
	"context"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/attendancetype"
	attendanceTypeRepo "github.com/kennykarnama/school-adminstration-api/domain/repository/attendancetype"
)

type Service interface {
	List(ctx context.Context) ([]*attendancetype.AttendanceType, error)
}

type svc struct {
	repo attendanceTypeRepo.Repository
}

func NewService(repo attendanceTypeRepo.Repository) *svc {
	return &svc{
		repo: repo,
	}
}

func (s *svc) List(ctx context.Context) ([]*attendancetype.AttendanceType, error) {
	return s.repo.List(ctx)
}
