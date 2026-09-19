package handler

import (
	"github.com/go-playground/validator/v10"
	userEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"github.com/kennykarnama/school-adminstration-api/domain/service/user"
	"github.com/kennykarnama/school-adminstration-api/util"
	"net/http"
	"time"
)

type UserHandler struct {
	svc          user.Service
	validate     *validator.Validate
	cookieSecure bool
	authEnabled  bool
}

func NewUserHandler(svc user.Service, validate *validator.Validate, cookieSecure, authEnabled bool) *UserHandler {
	return &UserHandler{
		svc:          svc,
		validate:     validate,
		cookieSecure: cookieSecure,
		authEnabled:  authEnabled,
	}
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := util.DecodeToStruct(r.Body, &req)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	if err := uh.validate.Struct(&req); err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	userSession, err := uh.svc.Login(r.Context(), req.AlternativeId, req.Password)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}
	x := time.Duration(userSession.Ttl) * time.Second
	expiredAt := userSession.CreatedAt.Add(x)
	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    userSession.Token,
		Path:     "/",
		Expires:  expiredAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   uh.cookieSecure,
	})
	resp := &LoginResponse{Token: userSession.Token}
	ResponseJson(w, resp, 201)
	return
}

func (uh *UserHandler) Validate(w http.ResponseWriter, r *http.Request) {
	if !uh.authEnabled {
		ResponseJson(w, TeacherProfileResponse{Name: "Administrator", Role: userEntity.RoleSchoolAdmin, School: &SchoolSummaryResponse{ID: "00000000-0000-4000-8000-000000000001", Name: "Existing School", Code: "legacy"}}, http.StatusOK)
		return
	}
	session, err := userEntity.NewUserSessionFromCtx(r.Context())
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
		return
	}
	teacher, err := uh.svc.Profile(r.Context(), session.UserId)
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	response := TeacherProfileResponse{
		ID:                 teacher.Id,
		AlternativeID:      teacher.AlternativeId,
		Name:               teacher.Name,
		Role:               teacher.Role,
		MustChangePassword: teacher.MustChangePassword,
	}
	if teacher.SchoolID != nil {
		response.School = &SchoolSummaryResponse{ID: *teacher.SchoolID, Name: teacher.SchoolName, Code: teacher.SchoolCode}
	}
	ResponseJson(w, response, http.StatusOK)
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if uh.authEnabled {
		session, err := userEntity.NewUserSessionFromCtx(r.Context())
		if err != nil {
			ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
			return
		}
		if err := uh.svc.Logout(r.Context(), session.Token); err != nil {
			ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
			return
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   uh.cookieSecure,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	principal, err := userEntity.NewPrincipalFromCtx(r.Context())
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
		return
	}
	var req ChangePasswordRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := uh.validate.Struct(req); err != nil || !validPassword(req.NewPassword) {
		ResponseJson(w, ErrorResponse{Message: "Kata sandi baru minimal 10 karakter serta mengandung huruf dan angka"}, http.StatusBadRequest)
		return
	}
	if err := uh.svc.ChangePassword(r.Context(), principal.UserID, req.CurrentPassword, req.NewPassword); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
