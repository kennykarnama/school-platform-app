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
	Profile(ctx context.Context, userID string) (*user2.Teacher, error)
	Logout(ctx context.Context, token string) error
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
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
	if !userData.Active {
		return nil, user2.ErrAccountInactive
	}
	if userData.SchoolID != nil && !userData.SchoolActive {
		return nil, user2.ErrAccountInactive
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

func (s *service) Profile(ctx context.Context, userID string) (*user2.Teacher, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *service) Logout(ctx context.Context, token string) error {
	return s.repo.DeleteUserSessionByToken(ctx, token)
}

func (s *service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	value, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if !util.PasswordMatch(currentPassword, value.Password) {
		return user2.ErrInvalidCredentials
	}
	hash, err := util.DefaultEncryptPassword(newPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, userID, hash, false)
}

func (s *service) RegisterTeachers(ctx context.Context, teachers []*user2.Teacher) error {
	for _, t := range teachers {
		hash, err := util.DefaultEncryptPassword(t.Password)
		if err != nil {
			return err
		}
		t.Password = hash
	}
	return s.repo.SaveTeachers(ctx, teachers)
}
