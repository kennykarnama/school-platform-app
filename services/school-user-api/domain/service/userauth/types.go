package userauth

import (
	"context"

	"github.com/kennykarnama/school-user-api/domain/models/auth"
	userEntity "github.com/kennykarnama/school-user-api/domain/models/user"
)

type LoginRequest struct {
	AlternativeID string `json:"alternativeID"`
	Password      string `json:"password"`
	DeviceID      string `json:"deviceID"`
	DeviceName    string `json:"deviceName"`
}

type ValidateTokenResult struct {
	User       *userEntity.User
	DeviceID   string
	DeviceName string
}

type Service interface {
	Login(ctx context.Context, req LoginRequest) (*auth.JWTMetadata, error)
	ValidateToken(ctx context.Context, token string) (*ValidateTokenResult, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*auth.JWTMetadata, error)
}
