package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	setupSvc "github.com/kennykarnama/school-adminstration-api/domain/service/setup"
)

const maxSetupRequestBytes = 2 << 20

type SetupHandler struct{ svc setupSvc.Service }

func NewSetupHandler(svc setupSvc.Service) *SetupHandler { return &SetupHandler{svc: svc} }

func (h *SetupHandler) Preview(w http.ResponseWriter, r *http.Request) {
	h.handle(w, r, false)
}

func (h *SetupHandler) Apply(w http.ResponseWriter, r *http.Request) {
	h.handle(w, r, true)
}

func (h *SetupHandler) handle(w http.ResponseWriter, r *http.Request, apply bool) {
	r.Body = http.MaxBytesReader(w, r.Body, maxSetupRequestBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var req setupSvc.Request
	if err := decoder.Decode(&req); err != nil {
		status := http.StatusBadRequest
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			status = http.StatusRequestEntityTooLarge
		}
		ResponseJson(w, ErrorResponse{Message: err.Error()}, status)
		return
	}
	if err := ensureJSONEOF(decoder); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	session, err := user.NewUserSessionFromCtx(r.Context())
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
		return
	}
	var result *setupSvc.Preview
	if apply {
		result, err = h.svc.Apply(r.Context(), req, session.UserId)
	} else {
		result, err = h.svc.Preview(r.Context(), req, session.UserId)
	}
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	status := http.StatusOK
	if apply && !result.Valid {
		status = http.StatusUnprocessableEntity
	}
	ResponseJson(w, result, status)
}

func ensureJSONEOF(decoder *json.Decoder) error {
	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("request body hanya boleh berisi satu objek JSON")
		}
		return err
	}
	return nil
}
