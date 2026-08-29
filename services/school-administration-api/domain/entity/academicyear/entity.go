package academicyear

import (
	"time"

	"gorm.io/gorm"
)

type AcademicYear struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Label     string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
