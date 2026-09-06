package academicyear

import (
	"context"

	entity "github.com/kennykarnama/school-adminstration-api/domain/entity/academicyear"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"github.com/kennykarnama/school-adminstration-api/domain/repository/academicyear"
)

type Service interface {
	List(ctx context.Context) ([]*entity.AcademicYear, error)
}

type service struct {
	repo academicyear.Repository
}

func NewService(repo academicyear.Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) List(ctx context.Context) ([]*entity.AcademicYear, error) {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.List(ctx, *principal)
}
