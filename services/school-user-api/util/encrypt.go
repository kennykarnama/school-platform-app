package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(plainPassword []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(plainPassword, bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("action=util.encrypt err=failed to encrypt. %v", err)
	}
	return string(hash), nil
}

func VerifyPassword(hashedPwd []byte, plainPwd []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	if err != nil {
		return false, fmt.Errorf("action=util.verifyPassword err=%v", err)
	}
	return true, nil
}
