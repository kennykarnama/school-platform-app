package setup

import (
	"context"
	"fmt"
	"strings"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/academicyear"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/attendancetype"
	"github.com/kennykarnama/school-adminstration-api/domain/entity/student"
	setupSvc "github.com/kennykarnama/school-adminstration-api/domain/service/setup"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type sqlRepository struct{ db *gorm.DB }

func NewSQLRepository(db *gorm.DB) *sqlRepository { return &sqlRepository{db: db} }

func (r *sqlRepository) Preview(ctx context.Context, req setupSvc.Request, teacherID string) (*setupSvc.Preview, error) {
	return buildPreview(r.db.WithContext(ctx), req, teacherID)
}

func (r *sqlRepository) Apply(ctx context.Context, req setupSvc.Request, teacherID string) (*setupSvc.Preview, error) {
	var result *setupSvc.Preview
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		preview, err := buildPreview(tx, req, teacherID)
		if err != nil {
			return err
		}
		result = preview
		if !preview.Valid {
			return nil
		}

		years, types, students, classes, err := loadState(tx)
		if err != nil {
			return err
		}
		yearByKey := uniqueYears(years)
		typeByKey := uniqueTypes(types)
		studentByKey := uniqueStudents(students)

		for _, input := range req.AcademicYears {
			key := fold(input.Label)
			if len(yearByKey[key]) == 0 {
				value := &academicyear.AcademicYear{ID: uuid.NewV4().String(), Label: input.Label}
				if err := tx.Create(value).Error; err != nil {
					return err
				}
				yearByKey[key] = []*academicyear.AcademicYear{value}
			}
		}
		for _, input := range req.AttendanceTypes {
			key := fold(input.Label)
			matches := typeByKey[key]
			if len(matches) == 0 {
				color := input.Color
				value := &attendancetype.AttendanceType{ID: uuid.NewV4().String(), Label: input.Label}
				if color != "" {
					value.Color = &color
				}
				if err := tx.Create(value).Error; err != nil {
					return err
				}
				typeByKey[key] = []*attendancetype.AttendanceType{value}
			} else if input.Color != "" {
				value := matches[0]
				if value.Color == nil || *value.Color != input.Color {
					if err := tx.Model(value).Update("color", input.Color).Error; err != nil {
						return err
					}
					color := input.Color
					value.Color = &color
				}
			}
		}

		classByKey := indexClasses(classes)
		for _, input := range req.Students {
			studentKey := fold(input.AlternativeID)
			var value *student.Student
			if len(studentByKey[studentKey]) == 0 {
				value = &student.Student{ID: uuid.NewV4().String(), AlternativeID: input.AlternativeID, Name: input.Name}
				if err := tx.Create(value).Error; err != nil {
					return err
				}
				studentByKey[studentKey] = []*student.Student{value}
			} else {
				value = studentByKey[studentKey][0]
				if value.Name != input.Name {
					if err := tx.Model(value).Update("name", input.Name).Error; err != nil {
						return err
					}
					value.Name = input.Name
				}
			}
			year := yearByKey[fold(input.AcademicYearLabel)][0]
			key := value.ID + "|" + year.ID
			matches := classByKey[key]
			if len(matches) == 0 {
				assignment := &student.StudentClass{ID: uuid.NewV4().String(), StudentID: value.ID, AcademicYearID: year.ID, ClassLabel: input.ClassLabel, UserID: teacherID}
				if err := tx.Create(assignment).Error; err != nil {
					return err
				}
				classByKey[key] = []*student.StudentClass{assignment}
			} else if matches[0].ClassLabel != input.ClassLabel {
				if err := tx.Model(matches[0]).Update("class_label", input.ClassLabel).Error; err != nil {
					return err
				}
				matches[0].ClassLabel = input.ClassLabel
			}
		}
		return nil
	})
	return result, err
}

func buildPreview(db *gorm.DB, req setupSvc.Request, teacherID string) (*setupSvc.Preview, error) {
	years, types, students, classes, err := loadState(db)
	if err != nil {
		return nil, err
	}
	yearByKey := uniqueYears(years)
	typeByKey := uniqueTypes(types)
	studentByKey := uniqueStudents(students)
	classByKey := indexClasses(classes)
	result := &setupSvc.Preview{Valid: true, Items: []setupSvc.ItemResult{}}

	for i, input := range req.AcademicYears {
		matches := yearByKey[fold(input.Label)]
		item := setupSvc.ItemResult{Index: i, Entity: "academicYear", Key: input.Label}
		switch len(matches) {
		case 0:
			item.Action = setupSvc.ActionCreate
		case 1:
			item.Action = setupSvc.ActionUnchanged
		default:
			item.Errors = []string{"Tahun ajaran ambigu karena duplikat sudah ada di database"}
		}
		add(result, item)
	}
	for i, input := range req.AttendanceTypes {
		matches := typeByKey[fold(input.Label)]
		item := setupSvc.ItemResult{Index: i, Entity: "attendanceType", Key: input.Label}
		switch len(matches) {
		case 0:
			item.Action = setupSvc.ActionCreate
		case 1:
			if input.Color != "" && (matches[0].Color == nil || *matches[0].Color != input.Color) {
				item.Action = setupSvc.ActionUpdate
			} else {
				item.Action = setupSvc.ActionUnchanged
			}
		default:
			item.Errors = []string{"Jenis kehadiran ambigu karena duplikat sudah ada di database"}
		}
		add(result, item)
	}
	for i, input := range req.Students {
		item := setupSvc.ItemResult{Index: i, Entity: "studentAssignment", Key: input.AlternativeID}
		yearMatches := yearByKey[fold(input.AcademicYearLabel)]
		if len(yearMatches) == 0 {
			// A year included in this same request will exist by apply time.
			for _, value := range req.AcademicYears {
				if fold(value.Label) == fold(input.AcademicYearLabel) {
					yearMatches = []*academicyear.AcademicYear{{ID: "new:" + fold(value.Label), Label: value.Label}}
					break
				}
			}
		}
		if len(yearMatches) == 0 {
			item.Errors = append(item.Errors, "Tahun ajaran tidak ditemukan")
		}
		if len(yearMatches) > 1 {
			item.Errors = append(item.Errors, "Tahun ajaran ambigu karena duplikat di database")
		}
		studentMatches := studentByKey[fold(input.AlternativeID)]
		if len(studentMatches) > 1 {
			item.Errors = append(item.Errors, "ID alternatif ambigu karena duplikat di database")
		}
		if len(item.Errors) == 0 && len(studentMatches) == 0 {
			item.Action = setupSvc.ActionCreate
		} else if len(item.Errors) == 0 {
			value := studentMatches[0]
			assignments := classByKey[value.ID+"|"+yearMatches[0].ID]
			if len(assignments) > 1 {
				item.Errors = append(item.Errors, "Penempatan siswa ambigu karena lebih dari satu data aktif")
			} else if len(assignments) == 1 && assignments[0].UserID != teacherID {
				item.Errors = append(item.Errors, "Penempatan siswa dimiliki guru lain")
			} else if value.Name != input.Name || len(assignments) == 0 || assignments[0].ClassLabel != input.ClassLabel {
				item.Action = setupSvc.ActionUpdate
			} else {
				item.Action = setupSvc.ActionUnchanged
			}
		}
		add(result, item)
	}
	return result, nil
}

func add(result *setupSvc.Preview, item setupSvc.ItemResult) {
	if len(item.Errors) > 0 {
		result.Valid = false
	} else {
		switch item.Action {
		case setupSvc.ActionCreate:
			result.Summary.Create++
		case setupSvc.ActionUpdate:
			result.Summary.Update++
		case setupSvc.ActionUnchanged:
			result.Summary.Unchanged++
		}
	}
	result.Items = append(result.Items, item)
}

func loadState(db *gorm.DB) ([]*academicyear.AcademicYear, []*attendancetype.AttendanceType, []*student.Student, []*student.StudentClass, error) {
	var years []*academicyear.AcademicYear
	var types []*attendancetype.AttendanceType
	var students []*student.Student
	var classes []*student.StudentClass
	if err := db.Find(&years).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	if err := db.Find(&types).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	if err := db.Find(&students).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	if err := db.Find(&classes).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	return years, types, students, classes, nil
}

func fold(value string) string { return strings.ToLower(strings.TrimSpace(value)) }
func uniqueYears(values []*academicyear.AcademicYear) map[string][]*academicyear.AcademicYear {
	result := map[string][]*academicyear.AcademicYear{}
	for _, value := range values {
		result[fold(value.Label)] = append(result[fold(value.Label)], value)
	}
	return result
}
func uniqueTypes(values []*attendancetype.AttendanceType) map[string][]*attendancetype.AttendanceType {
	result := map[string][]*attendancetype.AttendanceType{}
	for _, value := range values {
		result[fold(value.Label)] = append(result[fold(value.Label)], value)
	}
	return result
}
func uniqueStudents(values []*student.Student) map[string][]*student.Student {
	result := map[string][]*student.Student{}
	for _, value := range values {
		result[fold(value.AlternativeID)] = append(result[fold(value.AlternativeID)], value)
	}
	return result
}
func indexClasses(values []*student.StudentClass) map[string][]*student.StudentClass {
	result := map[string][]*student.StudentClass{}
	for _, value := range values {
		key := fmt.Sprintf("%s|%s", value.StudentID, value.AcademicYearID)
		result[key] = append(result[key], value)
	}
	return result
}
