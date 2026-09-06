package core

import (
	"context"
	"errors"
	"testing"

	studentEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	coreRepo "github.com/kennykarnama/school-adminstration-api/domain/repository/core"
)

type fakeStudentRepository struct {
	coreRepo.Repository
	updatedID       string
	updatedName     string
	updatePrincipal user.Principal
	listedOptions   studentEntity.StudentListOptions
	listedPrincipal user.Principal
	setActiveID     string
	setActiveValue  bool
	setActiveReason string
	setActivePrincipal user.Principal
	restoreClassID  string
	restorePrincipal user.Principal
}

func (r *fakeStudentRepository) UpdateStudentName(_ context.Context, studentID, name string, principal user.Principal) error {
	r.updatedID = studentID
	r.updatedName = name
	r.updatePrincipal = principal
	return nil
}

func (r *fakeStudentRepository) ListStudents(_ context.Context, options studentEntity.StudentListOptions, principal user.Principal) (*studentEntity.ManagementStudentPage, error) {
	r.listedOptions = options
	r.listedPrincipal = principal
	return &studentEntity.ManagementStudentPage{Items: []*studentEntity.ManagementStudent{}, Page: options.Page, PageSize: options.PageSize}, nil
}

func (r *fakeStudentRepository) SetStudentActive(_ context.Context, studentID string, active bool, reason string, principal user.Principal) error {
	r.setActiveID = studentID
	r.setActiveValue = active
	r.setActiveReason = reason
	r.setActivePrincipal = principal
	return nil
}

func (r *fakeStudentRepository) RestoreStudentClass(_ context.Context, studentClassID string, principal user.Principal) error {
	r.restoreClassID = studentClassID
	r.restorePrincipal = principal
	return nil
}

func TestUpdateStudentNameRequiresSchoolAdministrator(t *testing.T) {
	repo := &fakeStudentRepository{}
	service := NewService(repo, nil)
	ctx := user.Principal{UserID: "teacher-1", SchoolID: "school-1", Role: user.RoleTeacher}.NewCtx(context.Background())

	err := service.UpdateStudentName(ctx, "student-1", "New Name")
	if !errors.Is(err, user.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	if repo.updatedID != "" {
		t.Fatal("repository must not be called for a teacher")
	}
}

func TestUpdateStudentNameTrimsNameAndPassesTenantPrincipal(t *testing.T) {
	repo := &fakeStudentRepository{}
	service := NewService(repo, nil)
	ctx := user.Principal{UserID: "admin-1", SchoolID: "school-1", Role: user.RoleSchoolAdmin}.NewCtx(context.Background())

	if err := service.UpdateStudentName(ctx, "student-1", "  New Name  "); err != nil {
		t.Fatal(err)
	}
	if repo.updatedID != "student-1" || repo.updatedName != "New Name" || repo.updatePrincipal.SchoolID != "school-1" {
		t.Fatalf("unexpected repository call: id=%q name=%q principal=%+v", repo.updatedID, repo.updatedName, repo.updatePrincipal)
	}
}

func TestListStudentsPassesRoleAndFiltersToRepository(t *testing.T) {
	repo := &fakeStudentRepository{}
	service := NewService(repo, nil)
	ctx := user.Principal{UserID: "teacher-1", SchoolID: "school-1", Role: user.RoleTeacher}.NewCtx(context.Background())
	options := studentEntity.StudentListOptions{Query: "budi", AcademicYearID: "year-1", ClassLabel: "KELAS I A", Page: 2, PageSize: 50}

	if _, err := service.ListStudents(ctx, options); err != nil {
		t.Fatal(err)
	}
	if repo.listedOptions != options || repo.listedPrincipal.Role != user.RoleTeacher || repo.listedPrincipal.UserID != "teacher-1" {
		t.Fatalf("unexpected repository call: options=%+v principal=%+v", repo.listedOptions, repo.listedPrincipal)
	}
}

func TestSetStudentActiveRequiresSchoolAdministrator(t *testing.T) {
	repo := &fakeStudentRepository{}
	service := NewService(repo, nil)
	ctx := user.Principal{UserID: "teacher-1", SchoolID: "school-1", Role: user.RoleTeacher}.NewCtx(context.Background())

	err := service.SetStudentActive(ctx, "student-1", false, "alasan")
	if !errors.Is(err, user.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	if repo.setActiveID != "" {
		t.Fatal("repository must not be called for a teacher")
	}
}

func TestSetStudentActivePassesReasonAndPrincipal(t *testing.T) {
	repo := &fakeStudentRepository{}
	service := NewService(repo, nil)
	ctx := user.Principal{UserID: "admin-1", SchoolID: "school-1", Role: user.RoleSchoolAdmin}.NewCtx(context.Background())

	if err := service.SetStudentActive(ctx, "student-1", false, "  pindah sekolah  "); err != nil {
		t.Fatal(err)
	}
	if repo.setActiveID != "student-1" || repo.setActiveValue != false || repo.setActiveReason != "pindah sekolah" || repo.setActivePrincipal.SchoolID != "school-1" {
		t.Fatalf("unexpected repository call: id=%q active=%v reason=%q principal=%+v", repo.setActiveID, repo.setActiveValue, repo.setActiveReason, repo.setActivePrincipal)
	}
}

func TestRestoreStudentClassAllowsTeacherAndPassesPrincipal(t *testing.T) {
	repo := &fakeStudentRepository{}
	service := NewService(repo, nil)
	ctx := user.Principal{UserID: "teacher-1", SchoolID: "school-1", Role: user.RoleTeacher}.NewCtx(context.Background())

	if err := service.RestoreStudentClass(ctx, "class-1"); err != nil {
		t.Fatal(err)
	}
	if repo.restoreClassID != "class-1" || repo.restorePrincipal.Role != user.RoleTeacher || repo.restorePrincipal.UserID != "teacher-1" {
		t.Fatalf("unexpected repository call: id=%q principal=%+v", repo.restoreClassID, repo.restorePrincipal)
	}
}
