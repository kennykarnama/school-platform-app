package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"github.com/kennykarnama/school-adminstration-api/domain/service/administration"
	"github.com/kennykarnama/school-adminstration-api/util"
)

var schoolCodePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,30}[a-z0-9]$`)
var passwordLetterPattern = regexp.MustCompile(`[A-Za-z]`)
var passwordNumberPattern = regexp.MustCompile(`[0-9]`)

func validPassword(value string) bool {
	return len(value) >= 10 && passwordLetterPattern.MatchString(value) && passwordNumberPattern.MatchString(value)
}

type AdministrationHandler struct{ svc *administration.Service }

func NewAdministrationHandler(svc *administration.Service) *AdministrationHandler {
	return &AdministrationHandler{svc: svc}
}

func (h *AdministrationHandler) ListSchools(w http.ResponseWriter, r *http.Request) {
	values, err := h.svc.ListSchools(r.Context())
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	ResponseJson(w, map[string]interface{}{"items": values}, http.StatusOK)
}

func (h *AdministrationHandler) CreateSchool(w http.ResponseWriter, r *http.Request) {
	var req administration.CreateSchoolRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	req.Code = strings.ToLower(strings.TrimSpace(req.Code))
	if strings.TrimSpace(req.Name) == "" || !schoolCodePattern.MatchString(req.Code) || strings.TrimSpace(req.AdministratorName) == "" || strings.TrimSpace(req.AdministratorUsername) == "" || !validPassword(req.TemporaryPassword) {
		ResponseJson(w, ErrorResponse{Message: "Nama, kode sekolah, administrator, dan kata sandi sementara minimal 10 karakter dengan huruf dan angka wajib diisi"}, http.StatusBadRequest)
		return
	}
	value, err := h.svc.CreateSchool(r.Context(), req)
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusConflict)
		return
	}
	ResponseJson(w, value, http.StatusCreated)
}

func (h *AdministrationHandler) SetSchoolActive(w http.ResponseWriter, r *http.Request) {
	var req activeRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := h.svc.SetSchoolActive(r.Context(), mux.Vars(r)["id"], req.Active); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdministrationHandler) UpdateSchool(w http.ResponseWriter, r *http.Request) {
	var req schoolRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	req.Code = strings.ToLower(strings.TrimSpace(req.Code))
	if strings.TrimSpace(req.Name) == "" || !schoolCodePattern.MatchString(req.Code) {
		ResponseJson(w, ErrorResponse{Message: "Nama dan kode sekolah yang valid wajib diisi"}, http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateSchool(r.Context(), mux.Vars(r)["id"], req.Name, req.Code); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdministrationHandler) ListTeachers(w http.ResponseWriter, r *http.Request) {
	principal, _ := user.NewPrincipalFromCtx(r.Context())
	values, err := h.svc.ListTeachers(r.Context(), principal.SchoolID)
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	ResponseJson(w, map[string]interface{}{"items": values}, http.StatusOK)
}

func (h *AdministrationHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	principal, _ := user.NewPrincipalFromCtx(r.Context())
	var req administration.CreateTeacherRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.AlternativeID) == "" || strings.TrimSpace(req.Name) == "" || !validPassword(req.TemporaryPassword) {
		ResponseJson(w, ErrorResponse{Message: "Username, nama, dan kata sandi sementara minimal 10 karakter dengan huruf dan angka wajib diisi"}, http.StatusBadRequest)
		return
	}
	value, err := h.svc.CreateTeacher(r.Context(), principal.SchoolID, req)
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusConflict)
		return
	}
	ResponseJson(w, value, http.StatusCreated)
}

func (h *AdministrationHandler) SetTeacherActive(w http.ResponseWriter, r *http.Request) {
	principal, _ := user.NewPrincipalFromCtx(r.Context())
	var req activeRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := h.svc.SetTeacherActive(r.Context(), principal.SchoolID, mux.Vars(r)["id"], req.Active); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdministrationHandler) ResetTeacherPassword(w http.ResponseWriter, r *http.Request) {
	principal, _ := user.NewPrincipalFromCtx(r.Context())
	var req temporaryPasswordRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil || !validPassword(req.TemporaryPassword) {
		ResponseJson(w, ErrorResponse{Message: "Kata sandi sementara minimal 10 karakter serta mengandung huruf dan angka"}, http.StatusBadRequest)
		return
	}
	if err := h.svc.ResetTeacherPassword(r.Context(), principal.SchoolID, mux.Vars(r)["id"], req.TemporaryPassword); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdministrationHandler) ReplaceTeacherAccess(w http.ResponseWriter, r *http.Request) {
	principal, _ := user.NewPrincipalFromCtx(r.Context())
	var req accessRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := h.svc.ReplaceTeacherAccess(r.Context(), principal.SchoolID, mux.Vars(r)["id"], req.Items); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdministrationHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
	principal, _ := user.NewPrincipalFromCtx(r.Context())
	var req classRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateClass(r.Context(), principal.SchoolID, req.Label); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type activeRequest struct {
	Active bool `json:"active"`
}
type temporaryPasswordRequest struct {
	TemporaryPassword string `json:"temporaryPassword"`
}
type accessRequest struct {
	Items []administration.Access `json:"items"`
}
type classRequest struct {
	Label string `json:"label"`
}
type schoolRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
