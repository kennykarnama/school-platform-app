package class

import (
	"context"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/class"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	classRepo "github.com/kennykarnama/school-adminstration-api/domain/repository/class"
)

type Service interface {
	List(ctx context.Context) ([]*class.Class, error)
}

type svc struct {
	repo classRepo.Repository
}

func NewService(repo classRepo.Repository) *svc {
	return &svc{
		repo: repo,
	}
}

func (s *svc) List(ctx context.Context) ([]*class.Class, error) {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.List(ctx, *principal)
}
