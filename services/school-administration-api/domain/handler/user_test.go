package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	userEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/user"
)

type fakeUserService struct {
	profile     *userEntity.Teacher
	profileErr  error
	logoutToken string
	logoutErr   error
}

func (f *fakeUserService) Login(context.Context, string, string) (*userEntity.UserSession, error) {
	return nil, nil
}

func (f *fakeUserService) Validate(context.Context, string) (*userEntity.UserSession, error) {
	return nil, nil
}

func (f *fakeUserService) Profile(context.Context, string) (*userEntity.Teacher, error) {
	return f.profile, f.profileErr
}

func (f *fakeUserService) Logout(_ context.Context, token string) error {
	f.logoutToken = token
	return f.logoutErr
}

func (f *fakeUserService) RegisterTeachers(context.Context, []*userEntity.Teacher) error {
	return nil
}

func TestValidateReturnsAuthenticatedTeacherProfile(t *testing.T) {
	svc := &fakeUserService{profile: &userEntity.Teacher{
		Id:            "teacher-id",
		AlternativeId: "teacher.demo",
		Name:          "Demo Teacher",
	}}
	h := NewUserHandler(svc, validator.New(), true, true)
	session := &userEntity.UserSession{UserId: "teacher-id", Token: "session-token"}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/teacher/session/validate", nil)
	req = req.WithContext(session.NewCtx(req.Context()))
	res := httptest.NewRecorder()

	h.Validate(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	var profile TeacherProfileResponse
	if err := json.Unmarshal(res.Body.Bytes(), &profile); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if profile.ID != "teacher-id" || profile.AlternativeID != "teacher.demo" || profile.Name != "Demo Teacher" {
		t.Fatalf("unexpected profile: %+v", profile)
	}
}

func TestLogoutRevokesSessionAndExpiresCookie(t *testing.T) {
	svc := &fakeUserService{}
	h := NewUserHandler(svc, validator.New(), true, true)
	session := &userEntity.UserSession{UserId: "teacher-id", Token: "session-token"}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/teacher/logout", nil)
	req = req.WithContext(session.NewCtx(req.Context()))
	res := httptest.NewRecorder()

	h.Logout(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, res.Code)
	}
	if svc.logoutToken != "session-token" {
		t.Fatalf("expected session token to be revoked, got %q", svc.logoutToken)
	}
	cookies := res.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "session_token" || cookies[0].Value != "" || cookies[0].MaxAge != -1 {
		t.Fatalf("expected an expired session cookie, got %+v", cookies)
	}
	if !cookies[0].Secure || !cookies[0].HttpOnly {
		t.Fatalf("expected secure HTTP-only cookie, got %+v", cookies[0])
	}
}

func TestValidateReturnsLocalAdministratorWhenAuthDisabled(t *testing.T) {
	h := NewUserHandler(&fakeUserService{}, validator.New(), false, false)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/teacher/session/validate", nil)
	res := httptest.NewRecorder()

	h.Validate(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	var profile TeacherProfileResponse
	if err := json.Unmarshal(res.Body.Bytes(), &profile); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if profile.Name != "Administrator" {
		t.Fatalf("expected local administrator profile, got %+v", profile)
	}
}
