package userauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kennykarnama/school-user-api/config"
	"github.com/kennykarnama/school-user-api/domain/models/auth"
	userEntity "github.com/kennykarnama/school-user-api/domain/models/user"
	"github.com/kennykarnama/school-user-api/domain/repository/userauth"
	"github.com/kennykarnama/school-user-api/domain/service/user"
	"github.com/kennykarnama/school-user-api/util"
	uuid "github.com/satori/go.uuid"
)

type service struct {
	cfg          config.Config
	userService  user.Service
	userAuthRepo userauth.Repository
}

func NewService(cfg config.Config, userService user.Service, userAuthRepo userauth.Repository) *service {
	return &service{cfg: cfg, userService: userService, userAuthRepo: userAuthRepo}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*auth.JWTMetadata, error) {
	usr, err := s.userService.GetUserByAlternativeID(ctx, req.AlternativeID)
	if err != nil {
		return nil, err
	}
	correct, err := util.VerifyPassword(usr.PasswordAsByte(), []byte(req.Password))
	if err != nil {
		return nil, err
	}
	if !correct {
		return nil, ErrNotAuthorized
	}
	// generate JWT
	md, err := s.createToken(usr, req.DeviceID, req.DeviceName)
	if err != nil {
		return nil, err
	}

	// save session
	err = s.userAuthRepo.SaveJWTSession(ctx, usr.ID, md)
	if err != nil {
		return nil, err
	}

	return md, nil
}

func (s *service) ValidateToken(ctx context.Context, token string) (*ValidateTokenResult, error) {

	splitToken := strings.Split(token, " ")
	if splitToken[0] == "Bearer" {
		if len(splitToken) < 2 {
			return nil, fmt.Errorf("bearer without token")
		}

		token = splitToken[1]
	}

	var claims auth.TokenClaims
	if err := s.parseToken(&claims, token, s.cfg.AccessTokenCfg.Secret); err != nil {
		return nil, err
	}

	if claims.Expired() {
		return nil, ErrTokenExpired
	}

	exist, err := s.userAuthRepo.IsAccessTokenExist(ctx, claims.TokenID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrTokenRevoked
	}

	usr, err := s.userService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return &ValidateTokenResult{
		User:       usr,
		DeviceID:   claims.DeviceID,
		DeviceName: claims.DeviceName,
	}, nil
}

func (s *service) Logout(ctx context.Context, token string) error {

	var claim auth.TokenClaims

	if err := s.parseToken(&claim, token, s.cfg.AccessTokenCfg.Secret); err != nil {
		return err
	}

	return s.userAuthRepo.DeleteJWTSession(ctx, claim.TokenID)
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*auth.JWTMetadata, error) {
	splitToken := strings.Split(refreshToken, " ")
	if splitToken[0] == "Bearer" {
		if len(splitToken) < 2 {
			return nil, fmt.Errorf("bearer without token")
		}

		refreshToken = splitToken[1]
	}

	var claims auth.TokenClaims
	if err := s.parseToken(&claims, refreshToken, s.cfg.RefreshTokenCfg.Secret); err != nil {
		return nil, err
	}

	if claims.Expired() {
		return nil, ErrTokenExpired
	}

	exist, err := s.userAuthRepo.IsRefreshTokenExist(ctx, claims.TokenID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrTokenRevoked
	}

	err = s.userAuthRepo.DeleteJWTSession(ctx, claims.TokenID)
	if err != nil {
		return nil, err
	}

	// generate JWT
	md, err := s.createToken(&userEntity.User{
		ID: claims.UserID,
	}, claims.DeviceID, claims.DeviceName)
	if err != nil {
		return nil, err
	}

	// save session
	err = s.userAuthRepo.SaveJWTSession(ctx, claims.UserID, md)
	if err != nil {
		return nil, err
	}

	return md, nil
}

func (s *service) parseToken(claims *auth.TokenClaims, token string, secret string) error {
	_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		// TODO: logrus log
		fmt.Println(err)
		return ErrNotAuthorized
	}

	return nil
}

func (s *service) createToken(usr *userEntity.User, deviceID, deviceName string) (*auth.JWTMetadata, error) {
	tokenID := uuid.NewV4().String()

	refreshTokenClaims := &auth.TokenClaims{
		UserID:     usr.ID,
		TokenID:    tokenID,
		DeviceID:   deviceID,
		DeviceName: deviceName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(s.cfg.RefreshTokenCfg.Expiration).Unix(),
		},
	}
	refreshToken, err := s.createJwt(refreshTokenClaims, s.cfg.RefreshTokenCfg.Secret)
	if err != nil {
		return nil, err
	}

	accessTokenClaims := &auth.TokenClaims{
		UserID:     usr.ID,
		TokenID:    tokenID,
		DeviceID:   deviceID,
		DeviceName: deviceName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(s.cfg.AccessTokenCfg.Expiration).Unix(),
		},
	}
	accessToken, err := s.createJwt(accessTokenClaims, s.cfg.AccessTokenCfg.Secret)
	if err != nil {
		return nil, err
	}

	return &auth.JWTMetadata{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		TokenID:               tokenID,
		AccessTokenExpiresIn:  accessTokenClaims.ExpiresAt,
		RefreshTokenExpiresIn: refreshTokenClaims.ExpiresAt,
	}, nil
}

func (s *service) createJwt(tokenClaims *auth.TokenClaims, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("action=createJwt err=%v", err)
	}
	return token, nil
}
