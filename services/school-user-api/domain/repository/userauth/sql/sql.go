package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/kennykarnama/school-user-api/domain/models/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

type jwtSession struct {
	TokenID          string    `gorm:"column:token_id;type:uuid;primaryKey"`
	UserID           string    `gorm:"column:user_id;type:uuid;not null"`
	AccessExpiresAt  time.Time `gorm:"column:access_expires_at;not null"`
	RefreshExpiresAt time.Time `gorm:"column:refresh_expires_at;not null"`
	CreatedAt        time.Time `gorm:"column:created_at;not null"`
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) SaveJWTSession(ctx context.Context, userID string, metadata *auth.JWTMetadata) error {
	now := time.Now().UTC()
	session := &jwtSession{
		TokenID:          metadata.TokenID,
		UserID:           userID,
		AccessExpiresAt:  metadata.AccessTokenExpiresTime(),
		RefreshExpiresAt: metadata.RefreshTokenExpiresTime(),
		CreatedAt:        now,
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("jwt_session").Where("refresh_expires_at <= ?", now).Delete(&jwtSession{}).Error; err != nil {
			return err
		}
		return tx.Table("jwt_session").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "token_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"user_id", "access_expires_at", "refresh_expires_at", "created_at"}),
		}).Create(session).Error
	})
	if err != nil {
		return fmt.Errorf("action=repo.saveJWTSession userID=%v err=%w", userID, err)
	}
	return nil
}

func (r *repository) IsAccessTokenExist(ctx context.Context, tokenID string) (bool, error) {
	return r.exists(ctx, tokenID, "access_expires_at")
}

func (r *repository) IsRefreshTokenExist(ctx context.Context, tokenID string) (bool, error) {
	return r.exists(ctx, tokenID, "refresh_expires_at")
}

func (r *repository) exists(ctx context.Context, tokenID, expiryColumn string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("jwt_session").
		Where("token_id = ? AND "+expiryColumn+" > ?", tokenID, time.Now().UTC()).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("action=repo.jwtSessionExists tokenID=%v err=%w", tokenID, err)
	}
	return count > 0, nil
}

func (r *repository) DeleteJWTSession(ctx context.Context, tokenID string) error {
	if err := r.db.WithContext(ctx).Table("jwt_session").Where("token_id = ?", tokenID).Delete(&jwtSession{}).Error; err != nil {
		return fmt.Errorf("action=repo.deleteJWTSession tokenID=%v err=%w", tokenID, err)
	}
	return nil
}
