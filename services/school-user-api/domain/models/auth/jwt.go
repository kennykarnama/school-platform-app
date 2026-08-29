package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMetadata struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	TokenID               string `json:"refreshTokenID" redis:"refreshTokenID"`
	AccessTokenExpiresIn  int64  `json:"accessTokenExpiresIn" redis:"accessTokenExpiresIn"`
	RefreshTokenExpiresIn int64  `json:"refreshTokenExpiresAt" redis:"refreshTokenExpiresIn"`
}

func (jm *JWTMetadata) AccessTokenExpiresTime() time.Time {
	return time.Unix(jm.AccessTokenExpiresIn, 0)
}

func (jm *JWTMetadata) RefreshTokenExpiresTime() time.Time {
	return time.Unix(jm.RefreshTokenExpiresIn, 0)
}

type TokenClaims struct {
	UserID     string `json:"uid,omitempty"`
	TokenID    string `json:"tid,omitempty"`
	DeviceID   string `json:"did,omitempty"`
	DeviceName string `json:"dn,omitempty"`
	jwt.StandardClaims
}

type JWTSessionByID struct {
	UserID    string `redis:"userID"`
	ExpiresIn int64  `redis:"expiresIn"`
}

func (tc *TokenClaims) Expired() bool {
	now := time.Now().Unix()
	if tc.ExpiresAt <= now {
		return true
	}
	return false
}
