package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMetadata struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	TokenID               string `json:"refreshTokenID"`
	AccessTokenExpiresIn  int64  `json:"accessTokenExpiresIn"`
	RefreshTokenExpiresIn int64  `json:"refreshTokenExpiresAt"`
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

func (tc *TokenClaims) Expired() bool {
	now := time.Now().Unix()
	if tc.ExpiresAt <= now {
		return true
	}
	return false
}
