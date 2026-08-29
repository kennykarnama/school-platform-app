package userauth

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kennykarnama/school-user-api/domain/api/shared"
	"github.com/kennykarnama/school-user-api/domain/service/userauth"
	"github.com/kennykarnama/school-user-api/util"
)

type Handler struct {
	userAuthSvc userauth.Service
	ctx         context.Context
	validate    *validator.Validate
}

func NewHandler(ctx context.Context, validate *validator.Validate, userAuthService userauth.Service) *Handler {
	return &Handler{ctx: ctx, validate: validate, userAuthSvc: userAuthService}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := util.DecodeToStruct(r.Body, &req)
	if err != nil {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: err.Error(),
		}, shared.ErrorToHTTPStatus(err))
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	metadata, err := h.userAuthSvc.Login(h.ctx, userauth.LoginRequest{
		AlternativeID: req.AlternativeID,
		Password:      req.Password,
		DeviceID:      req.DeviceID,
		DeviceName:    req.DeviceName,
	})
	if err != nil {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: err.Error(),
		}, shared.ErrorToHTTPStatus(err))
		return
	}

	resp := NewLoginResponseFromJwt(metadata)

	shared.ResponseJson(w, resp, http.StatusOK)

	return
}

func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: "token required on header",
		}, http.StatusBadRequest)
		return
	}

	resp, err := h.userAuthSvc.ValidateToken(h.ctx, token)
	if err != nil {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: err.Error(),
		}, shared.ErrorToHTTPStatus(err))
		return
	}

	shared.ResponseJson(w, validateTokenResponse{
		UserID:     resp.User.ID,
		DeviceID:   resp.DeviceID,
		DeviceName: resp.DeviceName,
	}, http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: "token required on header",
		}, http.StatusBadRequest)
		return
	}

	err := h.userAuthSvc.Logout(h.ctx, token)
	if err != nil {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: err.Error(),
		}, shared.ErrorToHTTPStatus(err))
		return
	}

	shared.ResponseJson(w, shared.Empty{}, http.StatusOK)
	return
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: "refreshToken required on header",
		}, http.StatusBadRequest)
		return
	}

	metadata, err := h.userAuthSvc.RefreshToken(h.ctx, refreshToken)
	if err != nil {
		shared.ResponseJson(w, shared.ErrorResponse{
			Message: err.Error(),
		}, shared.ErrorToHTTPStatus(err))
		return
	}

	resp := NewLoginResponseFromJwt(metadata)

	shared.ResponseJson(w, resp, http.StatusCreated)

	return
}
