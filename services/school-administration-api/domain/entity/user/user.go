package user

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

var (
	ErrCtxSessionNotExist          = errors.New("user session not exist in context")
	ErrCtxSessionValueNotParseable = errors.New("user session value is not parseable")
)

type Teacher struct {
	Id            string         `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	AlternativeId string         `gorm:"alternative_id" json:"alternativeId"`
	Name          string         `json:"name"`
	Password      string         `json:"-"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt"`
}

type UserSession struct {
	Id        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId    string // reference to teacher.id
	Token     string
	Ttl       int
	CreatedAt time.Time
}

func (us *UserSession) NewCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, "userSessionCtx", *us)
}

func NewUserSessionFromCtx(ctx context.Context) (*UserSession, error) {
	v := ctx.Value("userSessionCtx")
	if v == nil {
		return nil, ErrCtxSessionNotExist
	}
	userSession, ok := v.(UserSession)
	if !ok {
		return nil, ErrCtxSessionValueNotParseable
	}
	return &userSession, nil
}
