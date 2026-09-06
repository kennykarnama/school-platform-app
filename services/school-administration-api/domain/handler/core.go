package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	statEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/stats"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	academicSvc "github.com/kennykarnama/school-adminstration-api/domain/service/academicyear"
	"github.com/kennykarnama/school-adminstration-api/domain/service/attendancetype"
	classSvc "github.com/kennykarnama/school-adminstration-api/domain/service/class"
	"github.com/kennykarnama/school-adminstration-api/domain/service/core"
	"github.com/kennykarnama/school-adminstration-api/util"
)

type Handler struct {
	coreSvc           core.Service
	academicYearSvc   academicSvc.Service
	classSvc          classSvc.Service
	validate          *validator.Validate
	attendanceTypeSvc attendancetype.Service
}

func NewHandler(coreSvc core.Service, academicYearSvc academicSvc.Service, classSvc classSvc.Service, validate *validator.Validate, attendanceTypeSvc attendancetype.Service) *Handler {
	return &Handler{
		coreSvc:           coreSvc,
		academicYearSvc:   academicYearSvc,
		classSvc:          classSvc,
		validate:          validate,
		attendanceTypeSvc: attendanceTypeSvc,
	}
}

func (h *Handler) RegisterStudent(w http.ResponseWriter, r *http.Request) {
	h.createStudent(w, r, false)
}

func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	h.createStudent(w, r, true)
}

func (h *Handler) createStudent(w http.ResponseWriter, r *http.Request, includeStudent bool) {
	var req RegisterStudentRequest

	err := util.DecodeToStruct(r.Body, &req)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.AlternativeID = strings.ToUpper(strings.TrimSpace(req.AlternativeID))
	req.AcademicYearID = strings.TrimSpace(req.AcademicYearID)
	req.ClassLabel = strings.ToUpper(strings.TrimSpace(req.ClassLabel))

	if err := h.validate.Struct(&req); err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	ts := time.Now().UTC()
	alternativeID := req.AlternativeID
	if alternativeID == "" {
		alternativeID, err = util.GenerateStudentAlternativeID()
		if err != nil {
			ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
	}

	newStudent := &student.Student{
		Name:          req.Name,
		AlternativeID: alternativeID,
		Active:        true,
		CreatedAt:     ts,
	}
	newStudentClass := &student.StudentClass{
		ClassLabel:     req.ClassLabel,
		AcademicYearID: req.AcademicYearID,
		CreatedAt:      ts,
	}

	err = h.coreSvc.RegisterStudent(r.Context(), newStudent, newStudentClass)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	if includeStudent {
		ResponseJson(w, map[string]interface{}{
			"id": newStudent.ID, "alternativeID": newStudent.AlternativeID, "name": newStudent.Name,
			"assignments": []student.ManagementAssignment{{StudentClassID: newStudentClass.ID, AcademicYearID: newStudentClass.AcademicYearID, ClassLabel: newStudentClass.ClassLabel}},
		}, http.StatusCreated)
		return
	}
	ResponseJson(w, Empty{}, http.StatusCreated)
}

func (h *Handler) ListStudents(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	page, err := positiveQueryInt(values.Get("page"), 1)
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: "page harus berupa bilangan positif"}, http.StatusBadRequest)
		return
	}
	pageSize, err := positiveQueryInt(values.Get("pageSize"), 50)
	if err != nil || pageSize > 100 {
		ResponseJson(w, ErrorResponse{Message: "pageSize harus antara 1 dan 100"}, http.StatusBadRequest)
		return
	}
	result, err := h.coreSvc.ListStudents(r.Context(), student.StudentListOptions{
		Query: strings.TrimSpace(values.Get("query")), AcademicYearID: strings.TrimSpace(values.Get("academicYearID")),
		ClassLabel: strings.ToUpper(strings.TrimSpace(values.Get("classLabel"))), Status: strings.TrimSpace(values.Get("status")),
		Page: page, PageSize: pageSize,
	})
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	ResponseJson(w, result, http.StatusOK)
}

func (h *Handler) UpdateStudentName(w http.ResponseWriter, r *http.Request) {
	var req UpdateStudentNameRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		ResponseJson(w, ErrorResponse{Message: "Nama siswa wajib diisi"}, http.StatusBadRequest)
		return
	}
	if err := h.coreSvc.UpdateStudentName(r.Context(), mux.Vars(r)["id"], req.Name); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func positiveQueryInt(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, strconv.ErrSyntax
	}
	return parsed, nil
}

func (h *Handler) SubmitAttendance(w http.ResponseWriter, r *http.Request) {
	var req SubmitAttendanceRequest
	err := util.DecodeToStruct(r.Body, &req)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	var items []*student.StudentAttendance
	ts := time.Now().UTC()
	for _, p := range req.Items {
		items = append(items, &student.StudentAttendance{
			StudentClassID:   p.StudentClassID,
			Attend:           p.Attend,
			CreatedAt:        ts,
			AttendanceDate:   p.AttendanceDate,
			AttendanceTypeID: p.AttendanceTypeID,
		})
	}
	err = h.coreSvc.SubmitAttendance(r.Context(), items)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	ResponseJson(w, Empty{}, http.StatusCreated)

	return
}

func (h *Handler) Students(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	academicYearID := urlParams.Get("academicYearID")
	classLabel := urlParams.Get("classLabel")
	attendanceDate := urlParams.Get("attendanceDate")

	if academicYearID == "" {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         "academicYearID is required",
		}, http.StatusBadRequest)
		return
	}

	if classLabel == "" {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         "classLabel is required",
		}, http.StatusBadRequest)
		return
	}

	if attendanceDate == "" {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         "attendanceDate is required",
		}, http.StatusBadRequest)
		return
	}

	items, err := h.coreSvc.ListAttendance(r.Context(), academicYearID, classLabel, attendanceDate)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	var respItems []*ListAttendanceItem

	for _, item := range items {
		respItems = append(respItems, &ListAttendanceItem{
			StudentID:           item.StudentID,
			StudentClassID:      item.StudentClassID,
			Name:                item.Name,
			Attend:              item.Attend,
			StudentAttendanceID: item.StudentAttendanceID,
			AttendanceDate:      item.AttendanceDate.Format("2006-01-02"),
			AttendanceTypeID:    item.AttendanceTypeID,
		})
	}

	ResponseJson(w, ListAttendanceResponse{
		Items: respItems,
	}, http.StatusOK)

	return
}

func (h *Handler) ListAcademicYear(w http.ResponseWriter, r *http.Request) {
	items, err := h.academicYearSvc.List(r.Context())
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	var respItems []*ListAcademicYearItem

	for _, item := range items {
		respItems = append(respItems, &ListAcademicYearItem{
			ID:        item.ID,
			Label:     item.Label,
			CreatedAt: item.CreatedAt,
		})
	}

	ResponseJson(w, ListAcademicYearResponse{
		Items: respItems,
	}, http.StatusOK)

	return
}

func (h *Handler) ListClasses(w http.ResponseWriter, r *http.Request) {
	items, err := h.classSvc.List(r.Context())
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	var respItems []*ListClassItem

	for _, item := range items {
		respItems = append(respItems, &ListClassItem{
			ID:    item.ID,
			Label: item.Label,
		})
	}

	ResponseJson(w, ListClassResponse{
		Items: respItems,
	}, http.StatusOK)

	return
}

func (h *Handler) ListAttendanceTypes(w http.ResponseWriter, r *http.Request) {
	items, err := h.attendanceTypeSvc.List(r.Context())
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}
	var respItems []*ListAttendanceTypeItem

	for _, item := range items {
		respItems = append(respItems, &ListAttendanceTypeItem{
			ID:    item.ID,
			Label: item.Label,
			Color: item.Color,
		})
	}

	ResponseJson(w, ListAttendanceType{
		Items: respItems,
	}, http.StatusOK)

	return
}

func (h *Handler) DeactivateStudentClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         "studentClassID is required",
		}, http.StatusBadRequest)
		return
	}
	var req DeactivateStudentClassRequest
	err := util.DecodeToStruct(r.Body, &req)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}
	if strings.TrimSpace(req.Reason) == "" {
		ResponseJson(w, ErrorResponse{Message: "Alasan penonaktifan wajib diisi"}, http.StatusBadRequest)
		return
	}

	err = h.coreSvc.DeactivateStudentClass(r.Context(), id, req.Reason)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SetStudentActive(w http.ResponseWriter, r *http.Request) {
	var req SetStudentActiveRequest
	if err := util.DecodeToStruct(r.Body, &req); err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	if !req.Active && strings.TrimSpace(req.Reason) == "" {
		ResponseJson(w, ErrorResponse{Message: "Alasan penonaktifan wajib diisi"}, http.StatusBadRequest)
		return
	}
	if err := h.coreSvc.SetStudentActive(r.Context(), mux.Vars(r)["id"], req.Active, req.Reason); err != nil {
		if errors.Is(err, student.ErrActivePlacementAlreadyExists) {
			ResponseJson(w, ErrorResponse{Message: "Siswa sudah memiliki penempatan aktif untuk tahun ajaran ini"}, http.StatusConflict)
			return
		}
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RestoreStudentClass(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		ResponseJson(w, ErrorResponse{Message: "studentClassID is required"}, http.StatusBadRequest)
		return
	}
	if err := h.coreSvc.RestoreStudentClass(r.Context(), id); err != nil {
		if errors.Is(err, student.ErrActivePlacementAlreadyExists) {
			ResponseJson(w, ErrorResponse{Message: "Siswa sudah memiliki penempatan aktif untuk tahun ajaran ini"}, http.StatusConflict)
			return
		}
		ResponseJson(w, ErrorResponse{Message: err.Error()}, ErrorToHTTPStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	academicYearID := urlParams.Get("academicYearID")
	classLabel := urlParams.Get("classLabel")
	from := urlParams.Get("from")
	to := urlParams.Get("to")
	statType := urlParams.Get("type")
	if academicYearID == "" {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         "academicYearID is required",
		}, http.StatusBadRequest)
		return
	}

	if classLabel == "" {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         "classLabel is required",
		}, http.StatusBadRequest)
		return
	}

	if from == "" {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         "attendanceDate.from is required",
		}, http.StatusBadRequest)
		return
	}

	if to == "" {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         "attendanceDate.to is required",
		}, http.StatusBadRequest)
		return
	}

	if statType == "" {
		statType = statEntity.DefaultStatType.String()
	}

	req := core.StatByRangeRequest{
		AcademicYearID: academicYearID,
		ClassLabel:     classLabel,
		From:           util.StrToPointer(from),
		To:             util.StrToPointer(to),
		Type:           statEntity.NewStatTypeFromStr(statType),
	}

	result, err := h.coreSvc.StatsByAttendanceType(r.Context(), req)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	var respItems []*StatsAttendanceItem

	for _, item := range result.Default {
		statItem := &StatsAttendanceItem{
			Name:            item.Name,
			AttendanceStats: []*AttendanceStatsItem{},
		}
		for _, stItem := range item.AttendanceStats {
			statItem.AttendanceStats = append(statItem.AttendanceStats, &AttendanceStatsItem{
				ID:     stItem.ID,
				Label:  stItem.Label,
				Counts: stItem.Count,
			})
		}
		statItem.Total = statEntity.AttendanceStatItems(item.AttendanceStats).Total()

		respItems = append(respItems, statItem)
	}

	classicalResp := &ClassicalStats{
		Items:         []*ClassicalStatItem{},
		StudentsTotal: result.Classical.StudentsTotal,
	}

	for _, classicalItem := range result.Classical.Items {
		classicalResp.Items = append(classicalResp.Items, &ClassicalStatItem{
			AttendanceStat: AttendanceStatsItem{
				ID:     classicalItem.AttendanceStat.ID,
				Label:  classicalItem.AttendanceStat.Label,
				Counts: classicalItem.AttendanceStat.Count,
			},
			AttendanceDate: classicalItem.AttendanceDate,
		})
	}

	apiResp := &StatsAttendanceResponse{
		Items:     respItems,
		Classical: classicalResp,
	}

	ResponseJson(w, apiResp, http.StatusOK)

	return
}

func (h *Handler) TransferStudentClass(w http.ResponseWriter, r *http.Request) {
	var req TransferStudentClassRequest
	err := util.DecodeToStruct(r.Body, &req)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		ResponseJson(w, ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	err = h.coreSvc.TransferStudentClass(r.Context(), req.SourceAcademicYearId, req.SourceClassLabel,
		req.DestinationAcademicYearId, req.DestinationClassLabel)
	if err != nil {
		ResponseJson(w, ErrorResponse{
			CustomErrorCode: "-",
			Message:         err.Error(),
		}, ErrorToHTTPStatus(err))
		return
	}
	ResponseJson(w, struct{}{}, http.StatusCreated)
}
