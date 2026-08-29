package userauth

import (
	authEntity "github.com/kennykarnama/school-user-api/domain/models/auth"
)

type loginRequest struct {
	AlternativeID string `json:"alternativeID" validate:"required"`
	Password      string `json:"password" validate:"required"`
	DeviceID      string `json:"deviceID" validate:"required"`
	DeviceName    string `json:"deviceName"`
}

type tokenResponse struct {
	Value     string `json:"value"`
	ExpiresIn int64  `json:"expiresIn"`
}

type loginResponse struct {
	AccessToken  *tokenResponse `json:"accessToken,omitempty"`
	RefreshToken *tokenResponse `json:"refreshToken,omitempty"`
}

type validateTokenResponse struct {
	UserID     string `json:"userID"`
	DeviceID   string `json:"deviceID"`
	DeviceName string `json:"deviceName"`
}

func NewLoginResponseFromJwt(metadata *authEntity.JWTMetadata) *loginResponse {
	return &loginResponse{
		AccessToken: &tokenResponse{
			Value:     metadata.AccessToken,
			ExpiresIn: metadata.AccessTokenExpiresIn,
		},
		RefreshToken: &tokenResponse{
			Value:     metadata.RefreshToken,
			ExpiresIn: metadata.RefreshTokenExpiresIn,
		},
	}
}
