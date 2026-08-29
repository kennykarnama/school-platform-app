package user

import (
	"context"
	"github.com/kennykarnama/school-adminstration-api/config"
	user2 "github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"github.com/kennykarnama/school-adminstration-api/domain/repository/user"
	"github.com/kennykarnama/school-adminstration-api/util"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Service interface {
	Login(ctx context.Context, alternativeId, password string) (*user2.UserSession, error)
	Validate(ctx context.Context, token string) (*user2.UserSession, error)
	RegisterTeachers(ctx context.Context, teachers []*user2.Teacher) error
}

type service struct {
	repo user.Repository
	cfg  config.Config
}

func NewService(repo user.Repository, cfg config.Config) *service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *service) Login(ctx context.Context, alternativeId, password string) (*user2.UserSession, error) {
	userData, err := s.repo.GetUserByAlternativeId(ctx, alternativeId)
	if err != nil {
		return nil, err
	}
	if !util.PasswordMatch(password, userData.Password) {
		return nil, user2.ErrInvalidCredentials
	}
	u1 := uuid.Must(uuid.NewV4(), err)
	if err != nil {
		return nil, err
	}
	userSession := &user2.UserSession{
		UserId: userData.Id,
		Token:  u1.String(),
		Ttl:    s.cfg.SessionTTL,
	}
	err = s.repo.SaveUserSession(ctx, userSession)
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (s *service) Validate(ctx context.Context, token string) (*user2.UserSession, error) {
	userSession, err := s.repo.GetUserSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	x := time.Duration(userSession.Ttl) * time.Second
	if userSession.CreatedAt.Add(x).Before(time.Now().UTC()) {
		return nil, user2.ErrSessionHasExpired
	}
	return userSession, nil
}

func (s *service) RegisterTeachers(ctx context.Context, teachers []*user2.Teacher) error {
	for _, t := range teachers {
		t.Password, _ = util.DefaultEncryptPassword(t.Password)
	}
	return s.repo.SaveTeachers(ctx, teachers)
}
