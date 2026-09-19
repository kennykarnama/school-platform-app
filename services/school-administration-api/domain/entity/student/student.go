package student

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID               string `gorm:"type:uuid;default:uuid_generate_v4()"`
	SchoolID         string
	Name             string
	AlternativeID    string
	Graduated        bool
	Active           bool   `gorm:"default:true"`
	DeactivateReason string
	CreatedAt        time.Time
	DeletedAt        gorm.DeletedAt
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

type ManagementAssignment struct {
	StudentClassID    string `json:"studentClassID"`
	AcademicYearID    string `json:"academicYearID"`
	AcademicYearLabel string `json:"academicYearLabel"`
	ClassLabel        string `json:"classLabel"`
	Active            bool   `json:"active"`
	DeactivateReason  string `json:"deactivateReason,omitempty"`
}

type ManagementStudent struct {
	ID               string                 `json:"id"`
	AlternativeID    string                 `json:"alternativeID"`
	Name             string                 `json:"name"`
	Active           bool                   `json:"active"`
	DeactivateReason string                 `json:"deactivateReason,omitempty"`
	Assignments      []ManagementAssignment `json:"assignments" gorm:"-"`
}

type ManagementStudentPage struct {
	Items    []*ManagementStudent `json:"items"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
	Total    int64                `json:"total"`
}

type StudentListOptions struct {
	Query          string
	AcademicYearID string
	ClassLabel     string
	Status         string
	Page           int
	PageSize       int
}
