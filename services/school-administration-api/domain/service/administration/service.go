package administration

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"github.com/kennykarnama/school-adminstration-api/util"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type School struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
}

type Teacher struct {
	ID                 string   `json:"id"`
	AlternativeID      string   `json:"alternativeId"`
	Name               string   `json:"name"`
	Role               string   `json:"role"`
	Active             bool     `json:"active"`
	MustChangePassword bool     `json:"mustChangePassword"`
	Access             []Access `gorm:"-" json:"access"`
}

type Access struct {
	AcademicYearID string `json:"academicYearId"`
	ClassLabel     string `json:"classLabel"`
}

type CreateTeacherRequest struct {
	AlternativeID     string `json:"alternativeId"`
	Name              string `json:"name"`
	TemporaryPassword string `json:"temporaryPassword"`
}

type CreateSchoolRequest struct {
	Name                  string `json:"name"`
	Code                  string `json:"code"`
	AdministratorName     string `json:"administratorName"`
	AdministratorUsername string `json:"administratorUsername"`
	TemporaryPassword     string `json:"temporaryPassword"`
}

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

func (s *Service) ListSchools(ctx context.Context) ([]School, error) {
	var values []School
	return values, s.db.WithContext(ctx).Table("school").Order("created_at ASC").Find(&values).Error
}

func (s *Service) CreateSchool(ctx context.Context, req CreateSchoolRequest) (*School, error) {
	hash, err := util.DefaultEncryptPassword(req.TemporaryPassword)
	if err != nil {
		return nil, err
	}
	value := &School{ID: uuid.NewV4().String(), Name: strings.TrimSpace(req.Name), Code: strings.ToLower(strings.TrimSpace(req.Code)), Active: true, CreatedAt: time.Now().UTC()}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("school").Create(value).Error; err != nil {
			return err
		}
		schoolID := value.ID
		admin := &user.Teacher{Id: uuid.NewV4().String(), SchoolID: &schoolID, AlternativeId: strings.TrimSpace(req.AdministratorUsername), Name: strings.TrimSpace(req.AdministratorName), Password: hash, Role: user.RoleSchoolAdmin, Active: true, MustChangePassword: true, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
		return tx.Create(admin).Error
	})
	return value, err
}

func (s *Service) SetSchoolActive(ctx context.Context, schoolID string, active bool) error {
	result := s.db.WithContext(ctx).Table("school").Where("id = ?", schoolID).Update("active", active)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *Service) UpdateSchool(ctx context.Context, schoolID, name, code string) error {
	result := s.db.WithContext(ctx).Table("school").Where("id = ?", schoolID).Updates(map[string]interface{}{
		"name": strings.TrimSpace(name), "code": strings.ToLower(strings.TrimSpace(code)), "updated_at": time.Now().UTC(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *Service) ListTeachers(ctx context.Context, schoolID string) ([]Teacher, error) {
	var values []Teacher
	if err := s.db.WithContext(ctx).Table("teacher").Where("school_id = ? AND role = ? AND deleted_at IS NULL", schoolID, user.RoleTeacher).Order("name ASC").Find(&values).Error; err != nil {
		return nil, err
	}
	for i := range values {
		if err := s.db.WithContext(ctx).Table("teacher_class_access").Select("academic_year_id, class_label").Where("school_id = ? AND teacher_id = ?", schoolID, values[i].ID).Order("academic_year_id, class_label").Find(&values[i].Access).Error; err != nil {
			return nil, err
		}
	}
	return values, nil
}

func (s *Service) CreateTeacher(ctx context.Context, schoolID string, req CreateTeacherRequest) (*Teacher, error) {
	hash, err := util.DefaultEncryptPassword(req.TemporaryPassword)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	value := &user.Teacher{Id: uuid.NewV4().String(), SchoolID: &schoolID, AlternativeId: strings.TrimSpace(req.AlternativeID), Name: strings.TrimSpace(req.Name), Password: hash, Role: user.RoleTeacher, Active: true, MustChangePassword: true, CreatedAt: now, UpdatedAt: now}
	if err := s.db.WithContext(ctx).Create(value).Error; err != nil {
		return nil, err
	}
	return &Teacher{ID: value.Id, AlternativeID: value.AlternativeId, Name: value.Name, Role: value.Role, Active: value.Active, MustChangePassword: true, Access: []Access{}}, nil
}

func (s *Service) SetTeacherActive(ctx context.Context, schoolID, teacherID string, active bool) error {
	result := s.db.WithContext(ctx).Table("teacher").Where("id = ? AND school_id = ? AND role = ?", teacherID, schoolID, user.RoleTeacher).Update("active", active)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return user.ErrForbidden
	}
	return nil
}

func (s *Service) ResetTeacherPassword(ctx context.Context, schoolID, teacherID, temporaryPassword string) error {
	hash, err := util.DefaultEncryptPassword(temporaryPassword)
	if err != nil {
		return err
	}
	result := s.db.WithContext(ctx).Table("teacher").Where("id = ? AND school_id = ? AND role = ?", teacherID, schoolID, user.RoleTeacher).Updates(map[string]interface{}{"password": hash, "must_change_password": true})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return user.ErrForbidden
	}
	return nil
}

func (s *Service) ReplaceTeacherAccess(ctx context.Context, schoolID, teacherID string, access []Access) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var teacherCount int64
		if err := tx.Table("teacher").Where("id = ? AND school_id = ? AND role = ?", teacherID, schoolID, user.RoleTeacher).Count(&teacherCount).Error; err != nil || teacherCount != 1 {
			if err != nil {
				return err
			}
			return user.ErrForbidden
		}
		for _, item := range access {
			var count int64
			if err := tx.Table("academic_year").Where("id = ? AND school_id = ? AND deleted_at IS NULL", item.AcademicYearID, schoolID).Count(&count).Error; err != nil || count != 1 {
				if err != nil {
					return err
				}
				return user.ErrForbidden
			}
			if err := tx.Table("school_class").Where("school_id = ? AND label = ? AND active = true", schoolID, strings.ToUpper(strings.TrimSpace(item.ClassLabel))).Count(&count).Error; err != nil || count != 1 {
				if err != nil {
					return err
				}
				return user.ErrForbidden
			}
		}
		if err := tx.Where("school_id = ? AND teacher_id = ?", schoolID, teacherID).Delete(&teacherClassAccess{}).Error; err != nil {
			return err
		}
		seen := map[string]bool{}
		for _, item := range access {
			item.ClassLabel = strings.ToUpper(strings.TrimSpace(item.ClassLabel))
			key := item.AcademicYearID + "|" + item.ClassLabel
			if seen[key] {
				continue
			}
			seen[key] = true
			row := teacherClassAccess{ID: uuid.NewV4().String(), SchoolID: schoolID, TeacherID: teacherID, AcademicYearID: item.AcademicYearID, ClassLabel: item.ClassLabel, CreatedAt: time.Now().UTC()}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) CreateClass(ctx context.Context, schoolID, label string) error {
	label = strings.ToUpper(strings.TrimSpace(label))
	if label == "" {
		return errors.New("class label is required")
	}
	return s.db.WithContext(ctx).Table("school_class").Create(map[string]interface{}{"id": uuid.NewV4().String(), "school_id": schoolID, "label": label, "active": true, "created_at": time.Now().UTC()}).Error
}

type teacherClassAccess struct {
	ID             string
	SchoolID       string
	TeacherID      string
	AcademicYearID string
	ClassLabel     string
	CreatedAt      time.Time
}

func (teacherClassAccess) TableName() string { return "teacher_class_access" }
