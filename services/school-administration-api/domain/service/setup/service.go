package setup

import (
	"context"
	"fmt"
	"strings"

	classEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/class"
)

type ClassService interface {
	List(ctx context.Context) ([]*classEntity.Class, error)
}

type Service interface {
	Preview(ctx context.Context, req Request, teacherID string) (*Preview, error)
	Apply(ctx context.Context, req Request, teacherID string) (*Preview, error)
}

type svc struct {
	repo     Repository
	classSvc ClassService
}

func NewService(repo Repository, classSvc ClassService) *svc {
	return &svc{repo: repo, classSvc: classSvc}
}

func (s *svc) Preview(ctx context.Context, req Request, teacherID string) (*Preview, error) {
	req, invalid, err := s.normalizeAndValidate(ctx, req)
	if err != nil || invalid != nil {
		return invalid, err
	}
	return s.repo.Preview(ctx, req, teacherID)
}

func (s *svc) Apply(ctx context.Context, req Request, teacherID string) (*Preview, error) {
	req, invalid, err := s.normalizeAndValidate(ctx, req)
	if err != nil || invalid != nil {
		return invalid, err
	}
	return s.repo.Apply(ctx, req, teacherID)
}

func (s *svc) normalizeAndValidate(ctx context.Context, req Request) (Request, *Preview, error) {
	classes, err := s.classSvc.List(ctx)
	if err != nil {
		return req, nil, err
	}
	validClasses := make(map[string]bool, len(classes))
	for _, item := range classes {
		validClasses[strings.ToLower(strings.TrimSpace(item.Label))] = true
	}

	result := &Preview{Valid: true, Items: []ItemResult{}}
	seenYears := map[string]bool{}
	for i := range req.AcademicYears {
		req.AcademicYears[i].Label = strings.TrimSpace(req.AcademicYears[i].Label)
		key := strings.ToLower(req.AcademicYears[i].Label)
		item := ItemResult{Index: i, Entity: "academicYear", Key: req.AcademicYears[i].Label}
		if key == "" {
			item.Errors = append(item.Errors, "Tahun ajaran wajib diisi")
		} else if seenYears[key] {
			item.Errors = append(item.Errors, "Tahun ajaran duplikat dalam data")
		}
		seenYears[key] = true
		result.Items = append(result.Items, item)
	}

	seenTypes := map[string]bool{}
	for i := range req.AttendanceTypes {
		req.AttendanceTypes[i].Label = strings.TrimSpace(req.AttendanceTypes[i].Label)
		key := strings.ToLower(req.AttendanceTypes[i].Label)
		item := ItemResult{Index: i, Entity: "attendanceType", Key: req.AttendanceTypes[i].Label}
		if key == "" {
			item.Errors = append(item.Errors, "Jenis kehadiran wajib diisi")
		} else if seenTypes[key] {
			item.Errors = append(item.Errors, "Jenis kehadiran duplikat dalam data")
		}
		seenTypes[key] = true
		result.Items = append(result.Items, item)
	}

	if len(req.Students) > MaxStudents {
		result.Items = append(result.Items, ItemResult{Entity: "students", Key: "limit", Errors: []string{fmt.Sprintf("Maksimal %d siswa per impor", MaxStudents)}})
	}
	seenAssignments := map[string]bool{}
	namesByAlternativeID := map[string]string{}
	for i := range req.Students {
		student := &req.Students[i]
		student.AlternativeID = strings.TrimSpace(student.AlternativeID)
		student.Name = strings.TrimSpace(student.Name)
		student.AcademicYearLabel = strings.TrimSpace(student.AcademicYearLabel)
		student.ClassLabel = strings.ToUpper(strings.TrimSpace(student.ClassLabel))
		altKey := strings.ToLower(student.AlternativeID)
		yearKey := strings.ToLower(student.AcademicYearLabel)
		assignmentKey := altKey + "|" + yearKey
		item := ItemResult{Index: i, Entity: "studentAssignment", Key: student.AlternativeID}
		if student.AlternativeID == "" {
			item.Errors = append(item.Errors, "ID alternatif wajib diisi")
		}
		if student.Name == "" {
			item.Errors = append(item.Errors, "Nama siswa wajib diisi")
		}
		if student.AcademicYearLabel == "" {
			item.Errors = append(item.Errors, "Tahun ajaran wajib diisi")
		}
		if student.ClassLabel == "" {
			item.Errors = append(item.Errors, "Kelas wajib diisi")
		} else if !validClasses[strings.ToLower(student.ClassLabel)] {
			item.Errors = append(item.Errors, "Kelas tidak dikenal")
		}
		if seenAssignments[assignmentKey] {
			item.Errors = append(item.Errors, "Penempatan siswa untuk tahun ajaran ini duplikat")
		}
		seenAssignments[assignmentKey] = true
		if prior, ok := namesByAlternativeID[altKey]; ok && !strings.EqualFold(prior, student.Name) {
			item.Errors = append(item.Errors, "ID alternatif yang sama memiliki nama berbeda")
		} else if altKey != "" {
			namesByAlternativeID[altKey] = student.Name
		}
		result.Items = append(result.Items, item)
	}

	for _, item := range result.Items {
		if len(item.Errors) > 0 {
			result.Valid = false
		}
	}
	if !result.Valid {
		return req, result, nil
	}
	return req, nil, nil
}
