package user

import "github.com/kennykarnama/school-user-api/domain/api/shared"

type registerUserRequest struct {
	AlternativeID string `json:"alternativeID" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Password      string `json:"password" validate:"required"`
}

type registerUserResponse struct {
	ID            string                `json:"id"`
	ErrorResponse *shared.ErrorResponse `json:"errorResponse,omitempty"`
}
