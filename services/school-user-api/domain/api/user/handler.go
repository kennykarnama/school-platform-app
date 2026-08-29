package user

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kennykarnama/school-user-api/domain/api/shared"
	userEntity "github.com/kennykarnama/school-user-api/domain/models/user"
	"github.com/kennykarnama/school-user-api/domain/service/user"
	"github.com/kennykarnama/school-user-api/util"
)

type Handler struct {
	userService user.Service
	validate    *validator.Validate
	ctx         context.Context
}

func NewHandler(ctx context.Context, userService user.Service, validate *validator.Validate) *Handler {
	return &Handler{
		userService: userService,
		validate:    validate,
		ctx:         ctx,
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

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

	newUser := &userEntity.User{
		AlternativeID: req.AlternativeID,
		Password:      req.Password,
		Name:          req.Name,
	}

	err = h.userService.RegisterUser(h.ctx, newUser)
	if err != nil {
		shared.ResponseJson(w, shared.ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, shared.ErrorToHTTPStatus(err))
		return
	}

	shared.ResponseJson(w, registerUserResponse{
		ID: newUser.ID,
	}, http.StatusCreated)

	return
}
