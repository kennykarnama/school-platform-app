package student

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID            string `gorm:"type:uuid;default:uuid_generate_v4()"`
	SchoolID      string
	Name          string
	AlternativeID string
	Graduated     bool
	CreatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

type StudentClass struct {
	ID               string `gorm:"type:uuid;default:uuid_generate_v4()"`
	SchoolID         string
	StudentID        string
	ClassLabel       string
	AcademicYearID   string
	CreatedAt        time.Time
	DeletedAt        gorm.DeletedAt
	DeactivateReason string
	UserID           string
}

type StudentAttendance struct {
	ID               string `gorm:"type:uuid;default:uuid_generate_v4()"`
	SchoolID         string
	StudentClassID   string
	Attend           bool
	CreatedAt        time.Time
	AttendanceDate   string
	AttendanceTypeID string
}

type Aggregate struct {
	StudentID           string
	StudentClassID      string
	Name                string
	Attend              bool
	StudentAttendanceID string
	AttendanceDate      time.Time
	AttendanceTypeID    string
}
