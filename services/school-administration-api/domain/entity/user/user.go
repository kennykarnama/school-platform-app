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

const (
	RolePlatformAdmin = "platform_admin"
	RoleSchoolAdmin   = "school_admin"
	RoleTeacher       = "teacher"
)

type Teacher struct {
	Id                 string         `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	AlternativeId      string         `gorm:"alternative_id" json:"alternativeId"`
	Name               string         `json:"name"`
	Password           string         `json:"-"`
	SchoolID           *string        `json:"schoolId"`
	Role               string         `json:"role"`
	Active             bool           `json:"active"`
	MustChangePassword bool           `json:"mustChangePassword"`
	SchoolName         string         `gorm:"column:school_name;->;-:migration" json:"schoolName,omitempty"`
	SchoolCode         string         `gorm:"column:school_code;->;-:migration" json:"schoolCode,omitempty"`
	SchoolActive       bool           `gorm:"column:school_active;->;-:migration" json:"-"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `json:"deletedAt"`
}

type Principal struct {
	UserID             string
	SchoolID           string
	Role               string
	MustChangePassword bool
}

type principalContextKey struct{}

func (p Principal) NewCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, principalContextKey{}, p)
}

func NewPrincipalFromCtx(ctx context.Context) (*Principal, error) {
	v := ctx.Value(principalContextKey{})
	if v == nil {
		return nil, ErrCtxSessionNotExist
	}
	p, ok := v.(Principal)
	if !ok {
		return nil, ErrCtxSessionValueNotParseable
	}
	return &p, nil
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
