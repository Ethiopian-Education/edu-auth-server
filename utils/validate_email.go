package utils

import (
	"regexp"
)

func ValidateEmail(email string) bool {
	r, _ := regexp.Compile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

	if !r.MatchString(email) {
		return false
	}
	return true
}
