package class

type Class struct {
	ID       string
	SchoolID string
	Label    string
	Active   bool
}

func (Class) TableName() string { return "school_class" }
