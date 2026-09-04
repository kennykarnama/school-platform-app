package user

import (
	"context"
	"fmt"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserByAlternativeId(ctx context.Context, alternativeId string) (*user.Teacher, error)
	GetUserByID(ctx context.Context, id string) (*user.Teacher, error)
	SaveUserSession(ctx context.Context, session *user.UserSession) error
	GetUserSessionByToken(ctx context.Context, token string) (*user.UserSession, error)
	DeleteUserSessionByToken(ctx context.Context, token string) error
	SaveTeachers(ctx context.Context, newTeachers []*user.Teacher) error
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) GetUserByAlternativeId(ctx context.Context, alternativeId string) (*user.Teacher, error) {
	var userData user.Teacher
	err := m.db.Where("alternative_id = ?", alternativeId).Find(&userData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("action=user.mysql.getUserByAlternativeId alternativeId=%s err=%v", alternativeId, err)
	}
	return &userData, nil
}

func (m *MySqlRepository) GetUserByID(ctx context.Context, id string) (*user.Teacher, error) {
	var userData user.Teacher
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&userData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("action=user.mysql.getUserByID id=%s err=%v", id, err)
	}
	return &userData, nil
}

func (m *MySqlRepository) SaveUserSession(ctx context.Context, session *user.UserSession) error {
	err := m.db.Create(session).Error
	if err != nil {
		return fmt.Errorf("action=user.mysql.saveUserSession err=%v", err)
	}
	return nil
}

func (m *MySqlRepository) GetUserSessionByToken(ctx context.Context, token string) (*user.UserSession, error) {
	var session user.UserSession
	err := m.db.Table("user_session").
		Where("token = ?", token).
		Order("created_at DESC").
		Limit(1).
		Find(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrSessionNotValid
		}
		return nil, fmt.Errorf("action=user.mysql.getUserSessionByToken err=%v", err)
	}
	return &session, nil
}

func (m *MySqlRepository) DeleteUserSessionByToken(ctx context.Context, token string) error {
	if err := m.db.WithContext(ctx).Where("token = ?", token).Delete(&user.UserSession{}).Error; err != nil {
		return fmt.Errorf("action=user.mysql.deleteUserSessionByToken err=%v", err)
	}
	return nil
}

func (m *MySqlRepository) SaveTeachers(ctx context.Context, newTeachers []*user.Teacher) error {
	err := m.db.Save(newTeachers).Error
	if err != nil {
		return fmt.Errorf("action=user.mysql.saveTeachers err=%v", err)
	}
	return nil
}
