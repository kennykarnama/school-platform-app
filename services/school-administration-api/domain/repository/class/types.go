package class

import (
	"context"
	"fmt"

	classEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/class"
)

type Repository interface {
	List(ctx context.Context) ([]*classEntity.Class, error)
}

type enumRepository struct{}

func NewEnumRepository() *enumRepository {
	return &enumRepository{}
}

func (r *enumRepository) List(ctx context.Context) ([]*classEntity.Class, error) {
	data := []*classEntity.Class{
		{
			ID:    "KELAS I",
			Label: "KELAS I",
		},
		{
			ID:    "KELAS II",
			Label: "KELAS II",
		},
		{
			ID:    "KELAS III",
			Label: "KELAS III",
		},
		{
			ID:    "KELAS IV",
			Label: "KELAS IV",
		},
		{
			ID:    "KELAS V",
			Label: "KELAS V",
		},
		{
			ID:    "KELAS VI",
			Label: "KELAS VI",
		},
	}
	classCategories := []string{
		"A", "B", "C",
	}
	var results []*classEntity.Class
	for _, item := range data {
		for _, ct := range classCategories {
			results = append(results, &classEntity.Class{
				ID:    fmt.Sprintf("%s %s", item.ID, ct),
				Label: fmt.Sprintf("%s %s", item.ID, ct),
			})
		}
	}
	return results, nil
}
