package academicyear

import (
	"time"

	"gorm.io/gorm"
)

type AcademicYear struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	SchoolID  string
	Label     string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
