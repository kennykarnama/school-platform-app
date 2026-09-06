package handler

import (
	"encoding/json"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"gorm.io/gorm"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, payload interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	rep, _ := json.Marshal(payload)
	w.Write(rep)
}

type Empty struct{}

type ErrorResponse struct {
	CustomErrorCode string `json:"customErrorCode,omitempty"`
	Message         string `json:"message,omitempty"`
}

func ErrorToHTTPStatus(err error) int {
	switch err {
	case user.ErrSessionHasExpired:
		return http.StatusUnauthorized
	case user.ErrSessionNotValid:
		return http.StatusUnauthorized
	case user.ErrInvalidCredentials:
		return http.StatusUnauthorized
	case user.ErrAccountInactive:
		return http.StatusUnauthorized
	case user.ErrForbidden, user.ErrPasswordChangeRequired:
		return http.StatusForbidden
	case student.ErrAlternativeIDAlreadyExists:
		return http.StatusConflict
	case student.ErrActivePlacementAlreadyExists:
		return http.StatusConflict
	case gorm.ErrRecordNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
