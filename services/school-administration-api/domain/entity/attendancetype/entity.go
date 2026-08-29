package attendancetype

import (
	"time"

	"gorm.io/gorm"
)

type AttendanceType struct {
	ID        string
	Label     string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type AttendanceTypes []*AttendanceType

func (ats AttendanceTypes) IDs() []string {
	var results []string
	for _, item := range ats {
		results = append(results, item.ID)
	}
	return results
}
