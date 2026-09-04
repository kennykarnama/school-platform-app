package setup

import "context"

const MaxStudents = 5000

type AcademicYearInput struct {
	Label string `json:"label"`
}

type AttendanceTypeInput struct {
	Label string `json:"label"`
	Color string `json:"color"`
}

type StudentInput struct {
	AlternativeID     string `json:"alternativeID"`
	Name              string `json:"name"`
	AcademicYearLabel string `json:"academicYearLabel"`
	ClassLabel        string `json:"classLabel"`
}

type Request struct {
	AcademicYears   []AcademicYearInput   `json:"academicYears"`
	AttendanceTypes []AttendanceTypeInput `json:"attendanceTypes"`
	Students        []StudentInput        `json:"students"`
}

type Action string

const (
	ActionCreate    Action = "create"
	ActionUpdate    Action = "update"
	ActionUnchanged Action = "unchanged"
)

type Summary struct {
	Create    int `json:"create"`
	Update    int `json:"update"`
	Unchanged int `json:"unchanged"`
}

type ItemResult struct {
	Index  int      `json:"index"`
	Entity string   `json:"entity"`
	Key    string   `json:"key"`
	Action Action   `json:"action,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

type Preview struct {
	Valid   bool         `json:"valid"`
	Summary Summary      `json:"summary"`
	Items   []ItemResult `json:"items"`
}

type Repository interface {
	Preview(ctx context.Context, req Request, teacherID string) (*Preview, error)
	Apply(ctx context.Context, req Request, teacherID string) (*Preview, error)
}

type ClassLister interface {
	List(ctx context.Context) ([]string, error)
}
