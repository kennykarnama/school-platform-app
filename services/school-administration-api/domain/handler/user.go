package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/kennykarnama/school-adminstration-api/domain/service/user"
	"github.com/kennykarnama/school-adminstration-api/util"
	"net/http"
	"time"
)

type UserHandler struct {
	svc          user.Service
	validate     *validator.Validate
	cookieSecure bool
}

func NewUserHandler(svc user.Service, validate *validator.Validate, cookieSecure bool) *UserHandler {
	return &UserHandler{
		svc:          svc,
		validate:     validate,
		cookieSecure: cookieSecure,
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
	ResponseJson(w, struct {
	}{}, 200)
	return
}
