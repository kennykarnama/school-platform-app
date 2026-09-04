package handler

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	setupSvc "github.com/kennykarnama/school-adminstration-api/domain/service/setup"
)

func TestStudentTemplateEndpoint(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewSetupHandler(nil).StudentTemplate(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/setup/students/template", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", recorder.Code, recorder.Body.String())
	}
	rows, err := readStudentXLSX(bytes.NewReader(recorder.Body.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := studentInputsFromRows(rows); err != nil {
		t.Fatalf("downloaded template cannot be imported: %v", err)
	}
}

func TestImportStudentsEndpoint(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "students.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("alternativeID,name,academicYearLabel,classLabel\n,Budi,2026,KELAS I A\n"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/setup/students/import", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	NewSetupHandler(nil).ImportStudents(recorder, request)
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"name":"Budi"`) {
		t.Fatalf("unexpected response: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestStudentInputsFromCSV(t *testing.T) {
	rows, err := readStudentCSV(strings.NewReader("alternativeID,name,academicYearLabel,classLabel\n,\"Putri, Ananda\",2026/2027,KELAS I A\n"))
	if err != nil {
		t.Fatal(err)
	}
	items, err := studentInputsFromRows(rows)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].AlternativeID != "" || items[0].Name != "Putri, Ananda" {
		t.Fatalf("unexpected imported students: %+v", items)
	}
}

func TestStudentTemplateRoundTrip(t *testing.T) {
	template, err := buildStudentTemplate()
	if err != nil {
		t.Fatal(err)
	}
	defer template.Close()
	buffer, err := template.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := readStudentXLSX(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	items, err := studentInputsFromRows(rows)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].AlternativeID != "" || items[0].Name != "Nama Siswa" {
		t.Fatalf("unexpected template contents: %+v", items)
	}
	styleID, err := template.GetCellStyle("Siswa", "A1")
	if err != nil || styleID == 0 {
		t.Fatalf("expected styled template header, style=%d err=%v", styleID, err)
	}
}

func TestStudentInputsRejectInvalidHeader(t *testing.T) {
	_, err := studentInputsFromRows([][]string{{"id", "name", "year", "class"}})
	if err == nil {
		t.Fatal("expected invalid header error")
	}
}

func TestStudentInputsRejectMoreThanLimit(t *testing.T) {
	rows := make([][]string, 1, setupSvc.MaxStudents+2)
	rows[0] = append([]string(nil), studentImportHeaders...)
	for index := 0; index <= setupSvc.MaxStudents; index++ {
		rows = append(rows, []string{"", fmt.Sprintf("Student %d", index), "2026", "KELAS I A"})
	}
	if _, err := studentInputsFromRows(rows); err == nil {
		t.Fatal("expected row limit error")
	}
}

func TestReadStudentXLSXRejectsInvalidWorkbook(t *testing.T) {
	if _, err := readStudentXLSX(strings.NewReader("not an xlsx file")); err == nil {
		t.Fatal("expected invalid workbook error")
	}
}
