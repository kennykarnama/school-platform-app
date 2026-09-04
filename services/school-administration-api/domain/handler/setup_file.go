package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	setupSvc "github.com/kennykarnama/school-adminstration-api/domain/service/setup"
	"github.com/xuri/excelize/v2"
)

const maxStudentImportBytes = 10 << 20

var studentImportHeaders = []string{"alternativeID", "name", "academicYearLabel", "classLabel"}

type StudentImportResponse struct {
	Items []setupSvc.StudentInput `json:"items"`
}

func (h *SetupHandler) StudentTemplate(w http.ResponseWriter, _ *http.Request) {
	workbook, err := buildStudentTemplate()
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	defer workbook.Close()

	buffer, err := workbook.WriteToBuffer()
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="template-data-awal-siswa.xlsx"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buffer.Bytes())
}

func (h *SetupHandler) ImportStudents(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxStudentImportBytes)
	file, header, err := r.FormFile("file")
	if err != nil {
		status := http.StatusBadRequest
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) || strings.Contains(err.Error(), "request body too large") {
			status = http.StatusRequestEntityTooLarge
		}
		ResponseJson(w, ErrorResponse{Message: "File Excel/CSV tidak dapat dibaca: " + err.Error()}, status)
		return
	}
	defer file.Close()

	var rows [][]string
	switch strings.ToLower(filepath.Ext(header.Filename)) {
	case ".csv":
		rows, err = readStudentCSV(file)
	case ".xlsx":
		rows, err = readStudentXLSX(file)
	default:
		err = errors.New("gunakan file Excel (.xlsx) atau CSV (.csv)")
	}
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	items, err := studentInputsFromRows(rows)
	if err != nil {
		ResponseJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	ResponseJson(w, StudentImportResponse{Items: items}, http.StatusOK)
}

func buildStudentTemplate() (*excelize.File, error) {
	workbook := excelize.NewFile()
	const sheet = "Siswa"
	if err := workbook.SetSheetName("Sheet1", sheet); err != nil {
		return nil, err
	}
	if err := workbook.SetSheetRow(sheet, "A1", &studentImportHeaders); err != nil {
		return nil, err
	}
	example := []interface{}{nil, "Nama Siswa", "2026/2027 - Semester 1", "KELAS I A"}
	if err := workbook.SetSheetRow(sheet, "A2", &example); err != nil {
		return nil, err
	}
	widths := []struct {
		column string
		width  float64
	}{{"A", 20}, {"B", 28}, {"C", 30}, {"D", 18}}
	for _, item := range widths {
		if err := workbook.SetColWidth(sheet, item.column, item.column, item.width); err != nil {
			return nil, err
		}
	}
	style, err := workbook.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"2563EB"}},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}
	if err := workbook.SetCellStyle(sheet, "A1", "D1", style); err != nil {
		return nil, err
	}
	if err := workbook.SetRowHeight(sheet, 1, 24); err != nil {
		return nil, err
	}
	if err := workbook.AutoFilter(sheet, "A1:D1", []excelize.AutoFilterOptions{}); err != nil {
		return nil, err
	}
	if err := workbook.SetPanes(sheet, &excelize.Panes{
		Freeze: true, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft",
	}); err != nil {
		return nil, err
	}
	return workbook, nil
}

func readStudentCSV(reader io.Reader) ([][]string, error) {
	parser := csv.NewReader(reader)
	parser.FieldsPerRecord = -1
	return parser.ReadAll()
}

func readStudentXLSX(reader io.Reader) ([][]string, error) {
	workbook, err := excelize.OpenReader(reader, excelize.Options{
		UnzipSizeLimit: 50 << 20, UnzipXMLSizeLimit: 10 << 20,
	})
	if err != nil {
		return nil, fmt.Errorf("file Excel tidak valid: %w", err)
	}
	defer workbook.Close()
	sheets := workbook.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("file Excel tidak memiliki lembar kerja")
	}
	return workbook.GetRows(sheets[0])
}

func studentInputsFromRows(rows [][]string) ([]setupSvc.StudentInput, error) {
	if len(rows) == 0 {
		return nil, errors.New("file tidak berisi data")
	}
	header := paddedStudentRow(rows[0])
	header[0] = strings.TrimPrefix(header[0], "\uFEFF")
	for index, expected := range studentImportHeaders {
		if strings.TrimSpace(header[index]) != expected {
			return nil, fmt.Errorf("header harus: %s", strings.Join(studentImportHeaders, ","))
		}
	}

	items := make([]setupSvc.StudentInput, 0, len(rows)-1)
	for _, source := range rows[1:] {
		row := paddedStudentRow(source)
		for index := range row {
			row[index] = strings.TrimSpace(row[index])
		}
		if row[0] == "" && row[1] == "" && row[2] == "" && row[3] == "" {
			continue
		}
		items = append(items, setupSvc.StudentInput{
			AlternativeID: row[0], Name: row[1], AcademicYearLabel: row[2], ClassLabel: row[3],
		})
		if len(items) > setupSvc.MaxStudents {
			return nil, fmt.Errorf("file melebihi %d baris siswa", setupSvc.MaxStudents)
		}
	}
	return items, nil
}

func paddedStudentRow(source []string) []string {
	row := make([]string, len(studentImportHeaders))
	copy(row, source)
	return row
}
