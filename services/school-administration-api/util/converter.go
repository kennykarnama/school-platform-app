package util

func StrToPointer(s string) *string {
	return &s
}

func PointerToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
