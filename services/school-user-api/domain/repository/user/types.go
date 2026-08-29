package user

import (
	"context"

	"github.com/kennykarnama/school-user-api/domain/models/user"
)

type Repository interface {
	RegisterUser(ctx context.Context, user *user.User) error
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	GetUserByAlternativeID(ctx context.Context, alternativeID string) (*user.User, error)
	UpdateUserNonEmpty(ctx context.Context, id string, user *user.User) error
}
