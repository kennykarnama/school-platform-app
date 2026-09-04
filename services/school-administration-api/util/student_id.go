package util

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// GenerateStudentAlternativeID returns a readable, cryptographically random
// identifier suitable for use as a student's public alternative ID.
func GenerateStudentAlternativeID() (string, error) {
	random := make([]byte, 8)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	return "STU-" + strings.ToUpper(hex.EncodeToString(random)), nil
}
