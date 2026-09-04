package util

import (
	"regexp"
	"testing"
)

func TestGenerateStudentAlternativeID(t *testing.T) {
	first, err := GenerateStudentAlternativeID()
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateStudentAlternativeID()
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^STU-[0-9A-F]{16}$`).MatchString(first) {
		t.Fatalf("unexpected generated ID: %s", first)
	}
	if first == second {
		t.Fatalf("expected distinct generated IDs, got %s", first)
	}
}
