package setup

import (
	"context"
	"testing"

	classEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/class"
)

type fakeRepository struct {
	request Request
	applied bool
}

func (r *fakeRepository) Preview(_ context.Context, req Request, _ string) (*Preview, error) {
	r.request = req
	return &Preview{Valid: true}, nil
}

func (r *fakeRepository) Apply(_ context.Context, req Request, _ string) (*Preview, error) {
	r.request = req
	r.applied = true
	return &Preview{Valid: true}, nil
}

type fakeClassService struct{}

func (fakeClassService) List(context.Context) ([]*classEntity.Class, error) {
	return []*classEntity.Class{{ID: "KELAS I A", Label: "KELAS I A"}}, nil
}

func TestPreviewNormalizesBeforeRepository(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, fakeClassService{})
	result, err := service.Preview(context.Background(), Request{
		AcademicYears:   []AcademicYearInput{{Label: " 2026/2027 "}},
		AttendanceTypes: []AttendanceTypeInput{{Label: " Sakit "}},
		Students:        []StudentInput{{AlternativeID: " S-1 ", Name: " Budi ", AcademicYearLabel: " 2026/2027 ", ClassLabel: " kelas i a "}},
	}, "teacher-1")
	if err != nil || !result.Valid {
		t.Fatalf("expected valid preview, result=%+v err=%v", result, err)
	}
	student := repo.request.Students[0]
	if student.AlternativeID != "S-1" || student.Name != "Budi" || student.ClassLabel != "KELAS I A" {
		t.Fatalf("request was not normalized: %+v", student)
	}
}

func TestPreviewRejectsInvalidAttendanceColor(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, fakeClassService{})
	result, err := service.Preview(context.Background(), Request{AttendanceTypes: []AttendanceTypeInput{{Label: "Sakit", Color: "red"}}}, "teacher-1")
	if err != nil || result.Valid {
		t.Fatalf("expected invalid color: %+v err=%v", result, err)
	}
	if repo.request.AttendanceTypes != nil {
		t.Fatal("repository must not run for invalid color")
	}
}

func TestPreviewGeneratesAlternativeIDsAndAllowsDuplicateNames(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, fakeClassService{})
	result, err := service.Preview(context.Background(), Request{Students: []StudentInput{
		{Name: "Budi", AcademicYearLabel: "2026", ClassLabel: "KELAS I A"},
		{Name: "Budi", AcademicYearLabel: "2026", ClassLabel: "KELAS I A"},
	}}, "teacher-1")
	if err != nil || !result.Valid {
		t.Fatalf("expected duplicate names with generated IDs to be valid: %+v err=%v", result, err)
	}
	first := repo.request.Students[0].AlternativeID
	second := repo.request.Students[1].AlternativeID
	if first == "" || second == "" || first == second {
		t.Fatalf("expected distinct generated alternative IDs, got %q and %q", first, second)
	}
}

func TestPreviewNormalizesAlternativeIDToUppercase(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, fakeClassService{})
	result, err := service.Preview(context.Background(), Request{Students: []StudentInput{{
		AlternativeID: " stu-abc123 ", Name: "Budi", AcademicYearLabel: "2026", ClassLabel: "KELAS I A",
	}}}, "teacher-1")
	if err != nil || !result.Valid {
		t.Fatalf("expected valid preview: %+v err=%v", result, err)
	}
	if got := repo.request.Students[0].AlternativeID; got != "STU-ABC123" {
		t.Fatalf("expected normalized ID, got %q", got)
	}
}

func TestPreviewRejectsDuplicateAssignmentAndConflictingName(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, fakeClassService{})
	result, err := service.Preview(context.Background(), Request{Students: []StudentInput{
		{AlternativeID: "S-1", Name: "Budi", AcademicYearLabel: "2026", ClassLabel: "KELAS I A"},
		{AlternativeID: "s-1", Name: "Andi", AcademicYearLabel: "2026", ClassLabel: "KELAS I A"},
	}}, "teacher-1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Valid || len(result.Items[1].Errors) < 2 {
		t.Fatalf("expected duplicate and conflicting-name errors: %+v", result)
	}
	if repo.request.Students != nil {
		t.Fatal("repository must not run for statically invalid input")
	}
}

func TestApplyRejectsUnknownClassWithoutWriting(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, fakeClassService{})
	result, err := service.Apply(context.Background(), Request{Students: []StudentInput{{
		AlternativeID: "S-1", Name: "Budi", AcademicYearLabel: "2026", ClassLabel: "KELAS X",
	}}}, "teacher-1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Valid || repo.applied {
		t.Fatalf("invalid class should not be applied: %+v", result)
	}
}
