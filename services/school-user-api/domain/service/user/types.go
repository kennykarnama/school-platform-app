package user

import (
	"context"

	"github.com/kennykarnama/school-user-api/domain/models/user"
)

type Service interface {
	RegisterUser(ctx context.Context, newUser *user.User) error
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	GetUserByAlternativeID(ctx context.Context, alternativeID string) (*user.User, error)
}
