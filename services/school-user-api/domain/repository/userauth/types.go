package userauth

import (
	"context"

	"github.com/kennykarnama/school-user-api/domain/models/auth"
)

type Repository interface {
	SaveJWTSession(ctx context.Context, userID string, metadata *auth.JWTMetadata) error
	IsAccessTokenExist(ctx context.Context, token string) (bool, error)
	IsRefreshTokenExist(ctx context.Context, token string) (bool, error)
	DeleteJWTSession(ctx context.Context, token string) error
}
