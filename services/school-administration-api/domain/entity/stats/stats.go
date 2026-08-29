package stats

type Stats struct {
	Name            string
	AttendanceStats []*AttendanceStatItem
	Total           int
}

type AttendanceStatItem struct {
	ID    string
	Label string
	Count int
}

type AttendanceStatItems []*AttendanceStatItem

func (as AttendanceStatItems) Total() int {
	var total int
	for _, item := range as {
		total += item.Count
	}
	return total
}

type ClassicalStats struct {
	Items         []*ClassicalStatItem
	StudentsTotal int
}

type ClassicalStatItem struct {
	AttendanceStat AttendanceStatItem
	AttendanceDate string
}

type StatType int

const (
	UnknownStatType StatType = iota
	DefaultStatType
	ClassicalStatType
)

var (
	StatTypeStr = map[string]StatType{
		"DEFAULT_STAT":   DefaultStatType,
		"CLASSICAL_STAT": ClassicalStatType,
	}
	StatTypeAsStr = map[StatType]string{
		DefaultStatType:   "DEFAULT_STAT",
		ClassicalStatType: "CLASSICAL_STAT",
		UnknownStatType:   "UNKNOWN_STAT",
	}
)

func NewStatTypeFromStr(s string) StatType {
	st, ok := StatTypeStr[s]
	if !ok {
		return UnknownStatType
	}
	return st
}

func (st StatType) String() string {
	return StatTypeAsStr[st]
}
