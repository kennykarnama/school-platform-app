package user

import (
	"context"
	"errors"

	userEntity "github.com/kennykarnama/school-user-api/domain/models/user"
	userRepo "github.com/kennykarnama/school-user-api/domain/repository/user"
	"github.com/kennykarnama/school-user-api/util"
)

type service struct {
	repo userRepo.Repository
}

func NewService(repo userRepo.Repository) *service {
	return &service{repo: repo}
}

func (s *service) RegisterUser(ctx context.Context, user *userEntity.User) error {
	hashedPwd, err := util.Encrypt(user.PasswordAsByte())
	if err != nil {
		return err
	}
	user.Password = hashedPwd
	err = s.repo.RegisterUser(ctx, user)
	if err != nil {
		if errors.Is(err, userRepo.ErrDuplicateEntry) {
			return ErrUserAlreadyExist
		}
		return err
	}
	return nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*userEntity.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) GetUserByAlternativeID(ctx context.Context, alternativeID string) (*userEntity.User, error) {
	return s.repo.GetUserByAlternativeID(ctx, alternativeID)
}
