package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string `gorm:"type:uuid;default:uuid_generate_v4()"`
	AlternativeID string
	Password      string
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (u *User) PasswordAsByte() []byte {
	return []byte(u.Password)
}
