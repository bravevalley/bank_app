package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	password := []byte(pwd)

	newPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("Can't Hash Password")
	}

	return string(newPassword), nil
}

func VerifyPassword(inputPwd, hashedPwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(inputPwd))
}
