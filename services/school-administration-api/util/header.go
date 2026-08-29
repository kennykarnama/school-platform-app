package util

import (
	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"strings"
)

func ExtractBearerToken(headerVal string) (string, error) {
	tokenString := headerVal
	if len(tokenString) == 0 {
		return "", user.ErrInvalidCredentials
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	return tokenString, nil
}
