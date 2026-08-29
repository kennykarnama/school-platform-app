package shared

import (
	"errors"
	"net/http"

	"github.com/kennykarnama/school-user-api/domain/service/user"
	"github.com/kennykarnama/school-user-api/domain/service/userauth"
)

type ErrorResponse struct {
	CustomErrorCode string `json:"customErrorCode,omitempty"`
	Message         string `json:"message,omitempty"`
}

func ErrorToHTTPStatus(err error) int {
	if errors.Is(userauth.ErrNotAuthorized, err) || errors.Is(userauth.ErrTokenExpired, err) || errors.Is(userauth.ErrTokenRevoked, err) {
		return http.StatusUnauthorized
	}
	if errors.Is(user.ErrUserNotExist, err) {
		return http.StatusNotFound
	}
	if errors.Is(user.ErrUserAlreadyExist, err) {
		return http.StatusSeeOther
	}
	return http.StatusInternalServerError
}
