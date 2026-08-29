package core

import statEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/stats"

type StatByRangeRequest struct {
	AcademicYearID string
	ClassLabel     string
	From           *string
	To             *string
	Type           statEntity.StatType
}

type StatByRangeResponse struct {
	Default   []*statEntity.Stats
	Classical *statEntity.ClassicalStats
}
