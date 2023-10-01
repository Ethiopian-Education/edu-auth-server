package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(plain_text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain_text), 14)

	return string(bytes), err
}