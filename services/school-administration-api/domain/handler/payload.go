package handler

import "time"

type RegisterStudentRequest struct {
	Name           string `json:"name" validate:"required"`
	AlternativeID  string `json:"alternativeID"`
	AcademicYearID string `json:"academicYearID" validate:"required"`
	ClassLabel     string `json:"classLabel" validate:"required"`
}

type AttendanceItem struct {
	StudentClassID   string `json:"studentClassID" validate:"required"`
	Attend           bool   `json:"attend" validate:"required"`
	AttendanceDate   string `json:"attendanceDate" validate:"required"`
	AttendanceTypeID string `json:"attendanceTypeID" validate:"required"`
}

type SubmitAttendanceRequest struct {
	Items []*AttendanceItem `json:"items"`
}

type ListAttendanceItem struct {
	StudentID           string `json:"studentID"`
	StudentClassID      string `json:"studentClassID"`
	Name                string `json:"name"`
	Attend              bool   `json:"attend"`
	StudentAttendanceID string `json:"studentAttendanceID"`
	AttendanceDate      string `json:"attendanceDate"`
	AttendanceTypeID    string `json:"attendanceTypeID"`
}

type ListAttendanceResponse struct {
	Items []*ListAttendanceItem `json:"items"`
}

type ListAcademicYearItem struct {
	ID        string    `json:"ID"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListAcademicYearResponse struct {
	Items []*ListAcademicYearItem `json:"items"`
}

type ListClassItem struct {
	ID    string `json:"ID"`
	Label string `json:"label"`
}

type ListClassResponse struct {
	Items []*ListClassItem `json:"items"`
}

type ListAttendanceTypeItem struct {
	ID    string  `json:"ID"`
	Label string  `json:"label"`
	Color *string `json:"color"`
}

type ListAttendanceType struct {
	Items []*ListAttendanceTypeItem `json:"items"`
}

type DeactivateStudentClassRequest struct {
	Reason string `json:"reason"`
}

type StatsAttendanceItem struct {
	Name            string                 `json:"name"`
	AttendanceStats []*AttendanceStatsItem `json:"statItems"`
	Total           int                    `json:"total"`
}

type AttendanceStatsItem struct {
	ID     string `json:"attendanceID"`
	Label  string `json:"attendanceLabel"`
	Counts int    `json:"counts"`
}

type StatsAttendanceResponse struct {
	Items     []*StatsAttendanceItem `json:"items"`
	Classical *ClassicalStats        `json:"classical"`
}

type ClassicalStats struct {
	Items         []*ClassicalStatItem `json:"items"`
	StudentsTotal int                  `json:"studentTotal"`
}

type ClassicalStatItem struct {
	AttendanceStat AttendanceStatsItem `json:"item"`
	AttendanceDate string              `json:"attendanceDate"`
}

type LoginRequest struct {
	AlternativeId string `json:"alternativeId" validate:"required"`
	Password      string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type TeacherProfileResponse struct {
	ID                 string                 `json:"id"`
	AlternativeID      string                 `json:"alternativeId"`
	Name               string                 `json:"name"`
	Role               string                 `json:"role"`
	MustChangePassword bool                   `json:"mustChangePassword"`
	School             *SchoolSummaryResponse `json:"school,omitempty"`
}

type SchoolSummaryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=10"`
}

type TransferStudentClassRequest struct {
	SourceAcademicYearId      string `json:"sourceAcademicYearId" validate:"required"`
	SourceClassLabel          string `json:"sourceClassLabel" validate:"required"`
	DestinationAcademicYearId string `json:"destinationAcademicYearId" validate:"required"`
	DestinationClassLabel     string `json:"destinationClassLabel" validate:"required"`
}
