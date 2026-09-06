package core

import (
	"context"
	"fmt"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"sort"
	"strings"
	"time"

	attendanceTypeEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/attendancetype"
	statEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/stats"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	coreRepo "github.com/kennykarnama/school-adminstration-api/domain/repository/core"
	"github.com/kennykarnama/school-adminstration-api/domain/service/attendancetype"
)

type Service interface {
	RegisterStudent(ctx context.Context, student *student.Student, class *student.StudentClass) error
	ListStudents(ctx context.Context, options student.StudentListOptions) (*student.ManagementStudentPage, error)
	UpdateStudentName(ctx context.Context, studentID, name string) error
	SubmitAttendance(ctx context.Context, data []*student.StudentAttendance) error
	ListAttendance(ctx context.Context, academicYearID string, classLabel string, attendanceDate string) ([]*student.Aggregate, error)
	DeactivateStudentClass(ctx context.Context, studentClassID string, reason string) error
	StatsByAttendanceType(ctx context.Context, req StatByRangeRequest) (*StatByRangeResponse, error)
	TransferStudentClass(ctx context.Context, sourceAcademicYear, sourceClass, destinationAcademicYear, destinationClass string) error
}

type svc struct {
	repo              coreRepo.Repository
	attendanceTypeSvc attendancetype.Service
}

func NewService(repo coreRepo.Repository, attendanceTypeSvc attendancetype.Service) *svc {
	return &svc{
		repo:              repo,
		attendanceTypeSvc: attendanceTypeSvc,
	}

}

func (s *svc) RegisterStudent(ctx context.Context, student *student.Student, class *student.StudentClass) error {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return err
	}
	return s.repo.RegisterStudent(ctx, student, class, *principal)
}

func (s *svc) ListStudents(ctx context.Context, options student.StudentListOptions) (*student.ManagementStudentPage, error) {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListStudents(ctx, options, *principal)
}

func (s *svc) UpdateStudentName(ctx context.Context, studentID, name string) error {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return err
	}
	if principal.Role != user.RoleSchoolAdmin {
		return user.ErrForbidden
	}
	return s.repo.UpdateStudentName(ctx, studentID, strings.TrimSpace(name), *principal)
}

func (s *svc) SubmitAttendance(ctx context.Context, data []*student.StudentAttendance) error {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return err
	}
	return s.repo.SubmitAttendance(ctx, data, *principal)
}

func (s *svc) ListAttendance(ctx context.Context, academicYearID string, classLabel string, attendanceDate string) ([]*student.Aggregate, error) {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListAttendance(ctx, academicYearID, classLabel, attendanceDate, *principal)
}

func (s *svc) DeactivateStudentClass(ctx context.Context, studentClassID string, reason string) error {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return err
	}
	return s.repo.DeactivateStudentClass(ctx, studentClassID, reason, *principal)
}

func (s *svc) StatsByAttendanceType(ctx context.Context, req StatByRangeRequest) (*StatByRangeResponse, error) {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	studentAttendances, err := s.repo.StudentAttendances(ctx, req.AcademicYearID, req.ClassLabel, req.From, req.To, *principal)
	if err != nil {
		return nil, err
	}
	studentClasses, err := s.repo.StudentClassesAggregate(ctx, req.AcademicYearID, req.ClassLabel, *principal)
	if err != nil {
		return nil, err
	}
	attendanceTypes, err := s.attendanceTypeSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	attendanceTypeIDs := attendanceTypeEntity.AttendanceTypes(attendanceTypes).IDs()
	var results []*statEntity.Stats
	for _, student := range studentClasses {
		var statItems []*statEntity.AttendanceStatItem
		for idx, attendanceTypeID := range attendanceTypeIDs {
			var counts int
			for _, stda := range studentAttendances {
				if stda.StudentClassID == student.StudentClassID && stda.AttendanceTypeID == attendanceTypeID {
					counts++
				}
			}
			statItems = append(statItems, &statEntity.AttendanceStatItem{
				ID:    attendanceTypeID,
				Label: attendanceTypes[idx].Label,
				Count: counts,
			})
		}
		statsByStudent := &statEntity.Stats{
			Name:            student.Name,
			AttendanceStats: statItems,
			Total:           statEntity.AttendanceStatItems(statItems).Total(),
		}
		results = append(results, statsByStudent)
	}

	classicalResults := &statEntity.ClassicalStats{
		Items:         []*statEntity.ClassicalStatItem{},
		StudentsTotal: len(studentClasses),
	}

	for _, attendanceType := range attendanceTypes {
		counters := make(map[string]int)
		attendanceDates := []string{}
		uniqueDatesTracker := make(map[string]bool)
		for _, stda := range studentAttendances {
			counters[fmt.Sprintf("%s:%s", stda.AttendanceTypeID, stda.AttendanceDate)]++
			if _, ok := uniqueDatesTracker[stda.AttendanceDate]; !ok {
				attendanceDates = append(attendanceDates, stda.AttendanceDate)
				uniqueDatesTracker[stda.AttendanceDate] = true
			}
		}
		for _, ad := range attendanceDates {
			classicalResults.Items = append(classicalResults.Items, &statEntity.ClassicalStatItem{
				AttendanceStat: statEntity.AttendanceStatItem{
					ID:    attendanceType.ID,
					Label: attendanceType.Label,
					Count: counters[fmt.Sprintf("%s:%s", attendanceType.ID, ad)],
				},
				AttendanceDate: ad,
			})
		}
	}
	sort.Slice(classicalResults.Items, func(i, j int) bool {
		d1, _ := time.Parse("2006-01-02T15:04:05Z", classicalResults.Items[i].AttendanceDate)
		d2, _ := time.Parse("2006-01-02T15:04:05Z", classicalResults.Items[j].AttendanceDate)
		return d1.Before(d2)
	})

	for _, ci := range classicalResults.Items {
		d1, err := time.Parse("2006-01-02T15:04:05Z", ci.AttendanceDate)
		if err != nil {
			return nil, err
		}
		ci.AttendanceDate = d1.Format("2006-01-02")
	}

	return &StatByRangeResponse{
		Default:   results,
		Classical: classicalResults,
	}, nil
}

func (s *svc) TransferStudentClass(ctx context.Context, sourceAcademicYear, sourceClass, destinationAcademicYear, destinationClass string) error {
	principal, err := user.NewPrincipalFromCtx(ctx)
	if err != nil {
		return err
	}
	if principal.Role == user.RoleTeacher {
		// Both sides must be explicitly assigned for teacher-initiated transfers.
		if _, err := s.repo.StudentClasses(ctx, destinationAcademicYear, destinationClass, *principal); err != nil {
			return err
		}
	}
	studentClasses, err := s.repo.StudentClasses(ctx, sourceAcademicYear, sourceClass, *principal)
	if err != nil {
		return err
	}
	for _, studentClass := range studentClasses {
		studentClass.AcademicYearID = destinationAcademicYear
		studentClass.ClassLabel = destinationClass
		studentClass.ID = ""
		studentClass.CreatedAt = time.Now().UTC().Add(4 * time.Second)
	}
	err = s.repo.SaveStudentClasses(ctx, studentClasses)
	if err != nil {
		return err
	}
	return nil
}
