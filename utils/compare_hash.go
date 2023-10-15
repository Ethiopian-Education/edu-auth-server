package utils

import "golang.org/x/crypto/bcrypt"

func CompareHashedPassword(hashed_text string, plain_text string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed_text), []byte(plain_text))

	return err == nil
}
