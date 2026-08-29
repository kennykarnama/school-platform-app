package sql

import (
	"context"
	"fmt"

	"github.com/kennykarnama/school-user-api/domain/models/user"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewMysqlRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) RegisterUser(ctx context.Context, user *user.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		return fmt.Errorf("action=repo.registeruser user.alternativeID=%v err:%v", user.AlternativeID, err)
	}
	return nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	var usr user.User
	err := r.db.Model(&user.User{}).Where("id = ?", id).Find(&usr).Error
	if err != nil {
		return nil, fmt.Errorf("action=repo.getUserByID id=%v err=%v", id, err)
	}
	// UUID empty will be considered as not exist user
	if usr.ID == "" {
		return nil, nil
	}
	return &usr, nil
}

func (r *repository) GetUserByAlternativeID(ctx context.Context, alternativeID string) (*user.User, error) {
	var usr user.User
	err := r.db.Model(&user.User{}).Where("alternative_id = ?", alternativeID).Find(&usr).Error
	if err != nil {
		return nil, fmt.Errorf("action=repo.getUserByAlternativeID id=%v err=%v", alternativeID, err)
	}
	// UUID empty will be considered as not exist user
	if usr.ID == "" {
		return nil, nil
	}
	return &usr, nil
}

func (r *repository) UpdateUserNonEmpty(ctx context.Context, id string, user *user.User) error {
	err := r.db.Model(user).Where("id = ?", id).Updates(*user).Error
	if err != nil {
		return fmt.Errorf("action=repo.updateUserNonEmpty id=%v err=%v", id, err)
	}
	return nil
}
