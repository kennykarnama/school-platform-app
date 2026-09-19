package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	coreSvc "github.com/kennykarnama/school-adminstration-api/domain/service/core"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
)

type fakeCoreService struct {
	setActiveID     string
	setActiveValue  bool
	setActiveReason string
	restoreClassID  string
	deactivateID    string
	deactivateReason string
	transferIDs     []string
	transferDestYear  string
	transferDestClass string
	transferCount   int
	transferErr     error
}

func (f *fakeCoreService) RegisterStudent(context.Context, *student.Student, *student.StudentClass) error { return nil }
func (f *fakeCoreService) ListStudents(context.Context, student.StudentListOptions) (*student.ManagementStudentPage, error) {
	return &student.ManagementStudentPage{}, nil
}
func (f *fakeCoreService) UpdateStudentName(context.Context, string, string) error { return nil }
func (f *fakeCoreService) SetStudentActive(_ context.Context, studentID string, active bool, reason string) error {
	f.setActiveID = studentID
	f.setActiveValue = active
	f.setActiveReason = reason
	return nil
}
func (f *fakeCoreService) RestoreStudentClass(_ context.Context, studentClassID string) error {
	f.restoreClassID = studentClassID
	return nil
}
func (f *fakeCoreService) SubmitAttendance(context.Context, []*student.StudentAttendance) error { return nil }
func (f *fakeCoreService) ListAttendance(context.Context, string, string, string) ([]*student.Aggregate, error) {
	return nil, nil
}
func (f *fakeCoreService) DeactivateStudentClass(_ context.Context, studentClassID string, reason string) error {
	f.deactivateID = studentClassID
	f.deactivateReason = reason
	return nil
}
func (f *fakeCoreService) StatsByAttendanceType(context.Context, coreSvc.StatByRangeRequest) (*coreSvc.StatByRangeResponse, error) { return nil, nil }
func (f *fakeCoreService) TransferStudentClass(context.Context, string, string, string, string) error { return nil }
func (f *fakeCoreService) TransferStudents(_ context.Context, studentClassIDs []string, destYear, destClass string) (int, error) {
	f.transferIDs = studentClassIDs
	f.transferDestYear = destYear
	f.transferDestClass = destClass
	return f.transferCount, f.transferErr
}

func newCoreHandler(svc *fakeCoreService) *Handler {
	return &Handler{
		coreSvc:  svc,
		validate: validator.New(),
	}
}

func TestSetStudentActiveRejectsEmptyReasonWhenDeactivating(t *testing.T) {
	svc := &fakeCoreService{}
	h := newCoreHandler(svc)
	body := strings.NewReader(`{"active":false,"reason":""}`)
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/students/student-1/status", body)
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "student-1"})
	w := httptest.NewRecorder()

	h.SetStudentActive(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Alasan penonaktifan wajib diisi" {
		t.Fatalf("unexpected message: %q", resp.Message)
	}
	if svc.setActiveID != "" {
		t.Fatal("service must not be called when reason is empty")
	}
}

func TestSetStudentActiveCallsServiceWithCorrectArgs(t *testing.T) {
	svc := &fakeCoreService{}
	h := newCoreHandler(svc)
	body := strings.NewReader(`{"active":false,"reason":"pindah sekolah"}`)
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/students/student-1/status", body)
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "student-1"})
	w := httptest.NewRecorder()

	h.SetStudentActive(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if svc.setActiveID != "student-1" || svc.setActiveValue != false || svc.setActiveReason != "pindah sekolah" {
		t.Fatalf("unexpected service call: id=%q active=%v reason=%q", svc.setActiveID, svc.setActiveValue, svc.setActiveReason)
	}
}

func TestSetStudentActiveAllowsEmptyReasonWhenActivating(t *testing.T) {
	svc := &fakeCoreService{}
	h := newCoreHandler(svc)
	body := strings.NewReader(`{"active":true,"reason":""}`)
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/students/student-1/status", body)
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "student-1"})
	w := httptest.NewRecorder()

	h.SetStudentActive(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if svc.setActiveID != "student-1" || svc.setActiveValue != true {
		t.Fatalf("unexpected service call: id=%q active=%v", svc.setActiveID, svc.setActiveValue)
	}
}

func TestDeactivateStudentClassReturns204AndRequiresReason(t *testing.T) {
	svc := &fakeCoreService{}
	h := newCoreHandler(svc)
	body := strings.NewReader(`{"reason":""}`)
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/student/class/class-1/deactivate", body)
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "class-1"})
	w := httptest.NewRecorder()

	h.DeactivateStudentClass(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty reason, got %d", w.Code)
	}
	if svc.deactivateID != "" {
		t.Fatal("service must not be called when reason is empty")
	}
}

func TestDeactivateStudentClassReturns204OnSuccess(t *testing.T) {
	svc := &fakeCoreService{}
	h := newCoreHandler(svc)
	body := strings.NewReader(`{"reason":"pindah kelas"}`)
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/student/class/class-1/deactivate", body)
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "class-1"})
	w := httptest.NewRecorder()

	h.DeactivateStudentClass(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if svc.deactivateID != "class-1" || svc.deactivateReason != "pindah kelas" {
		t.Fatalf("unexpected service call: id=%q reason=%q", svc.deactivateID, svc.deactivateReason)
	}
}

func TestRestoreStudentClassReturns204(t *testing.T) {
	svc := &fakeCoreService{}
	h := newCoreHandler(svc)
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/student/class/class-1/restore", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "class-1"})
	w := httptest.NewRecorder()

	h.RestoreStudentClass(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if svc.restoreClassID != "class-1" {
		t.Fatalf("unexpected service call: id=%q", svc.restoreClassID)
	}
}

func TestRestoreStudentClassReturns409OnClash(t *testing.T) {
	svc := &clashingCoreService{}
	h := &Handler{coreSvc: svc, validate: validator.New()}
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/student/class/class-1/restore", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "class-1"})
	w := httptest.NewRecorder()

	h.RestoreStudentClass(w, r)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Siswa sudah memiliki penempatan aktif untuk tahun ajaran ini" {
		t.Fatalf("unexpected message: %q", resp.Message)
	}
}

type clashingCoreService struct{ fakeCoreService }

func (f *clashingCoreService) RestoreStudentClass(context.Context, string) error {
	return student.ErrActivePlacementAlreadyExists
}

func TestSetStudentActiveReturns403ForTeacher(t *testing.T) {
	svc := &forbiddenCoreService{}
	h := &Handler{coreSvc: svc, validate: validator.New()}
	body := strings.NewReader(`{"active":false,"reason":"alasan"}`)
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/students/student-1/status", body)
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"id": "student-1"})
	w := httptest.NewRecorder()

	h.SetStudentActive(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

type forbiddenCoreService struct{ fakeCoreService }

func (f *forbiddenCoreService) SetStudentActive(context.Context, string, bool, string) error {
	return user.ErrForbidden
}

func TestTransferStudentsRejectsEmptyIDs(t *testing.T) {
	svc := &fakeCoreService{}
	h := newCoreHandler(svc)
	body := strings.NewReader(`{"studentClassIDs":[],"destinationAcademicYearId":"year-1","destinationClassLabel":"KELAS I A"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/student/transfer", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.TransferStudents(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if svc.transferIDs != nil {
		t.Fatal("service must not be called with empty IDs")
	}
}

func TestTransferStudentsReturnsCount(t *testing.T) {
	svc := &fakeCoreService{transferCount: 3}
	h := newCoreHandler(svc)
	body := strings.NewReader(`{"studentClassIDs":["id-1","id-2","id-3"],"destinationAcademicYearId":"year-1","destinationClassLabel":"KELAS I B"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/student/transfer", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.TransferStudents(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if len(svc.transferIDs) != 3 || svc.transferDestYear != "year-1" || svc.transferDestClass != "KELAS I B" {
		t.Fatalf("unexpected service call: ids=%v year=%q class=%q", svc.transferIDs, svc.transferDestYear, svc.transferDestClass)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["transferred"] != float64(3) {
		t.Fatalf("expected transferred=3, got %v", resp["transferred"])
	}
}

func TestTransferStudentsReturns409OnClash(t *testing.T) {
	svc := &fakeCoreService{transferErr: student.ErrActivePlacementAlreadyExists}
	h := newCoreHandler(svc)
	body := strings.NewReader(`{"studentClassIDs":["id-1"],"destinationAcademicYearId":"year-1","destinationClassLabel":"KELAS I A"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/student/transfer", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.TransferStudents(w, r)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}
