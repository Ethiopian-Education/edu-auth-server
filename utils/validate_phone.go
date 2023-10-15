package utils

import (
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

func ValidatePhone(phone string) (string, bool) {
	r, _ := regexp.Compile(`^(^\+251|^251|^0)?(9|7)\d{8}$`)

	if len(phone) < 9 || len(phone) > 13 {
		logrus.Info("Invalid length")
		return "", false
	}
	// Check the validity
	if !r.MatchString(phone) {
		return "", false
	}
	if strings.HasPrefix(phone, "9") {
		phone = strings.Replace(phone, "9", "+2519", 1)
	} else if strings.HasPrefix(phone, "251") {
		phone = strings.Replace(phone, "251", "+251", 1)
	} else if strings.HasPrefix(phone, "09") {
		phone = strings.Replace(phone, "09", "+2519", 1)
	} else if strings.HasPrefix(phone, "07") {
		phone = strings.Replace(phone, "07", "+2517", 1)
	} else if strings.HasPrefix(phone, "7") {
		phone = strings.Replace(phone, "7", "+2517", 1)
	}
	// logrus.Info("Phone -- ", phone)
	return phone, true
}
