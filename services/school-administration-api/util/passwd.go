package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func DefaultEncryptPassword(p string) (string, error) {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("action=util.defaultEncrypt err=%v", err)
	}
	return string(hashedPassword), nil
}

func PasswordMatch(target, hashed string) bool {
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(target))
	return err == nil
}
