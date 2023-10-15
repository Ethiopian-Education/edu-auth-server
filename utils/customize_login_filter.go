package utils

import (
	"fmt"
	"strings"
)

func NormalizeMultiLoginFilter(username string) string {
	phone, isValid := ValidatePhone(username)
	if isValid {
		filter_string := fmt.Sprintf(`phone_number:{_eq: "%s"}`, strings.TrimSpace(phone))
		return filter_string
	} else if ValidateEmail(username) {
		filter_string := fmt.Sprintf(`email:{_eq: "%s"}`, strings.TrimSpace(username))
		return filter_string
	} else {
		filter_string := fmt.Sprintf(`username:{_eq: "%s"}`, strings.TrimSpace(username))
		return filter_string
	}
}
