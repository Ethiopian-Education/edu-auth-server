package utils

import (
	"math/rand"
	"time"
)

func GenerateNumber() int {
	// Generate between a given range
	// digits := []rune{0,1,2,3,4,5,6,7,8,9}
	rand.Seed(time.Now().UnixNano())
	min := 4
	max := 7

	return rand.Intn(max-min) + min // range is min to max
}

func GenerateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())
	// chars that can be used in the charset
	digits_charset := "0123456789"
	otp := make([]byte, length)

	// Generate random chars from the digits_charset
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(digits_charset))
		otp[i] = digits_charset[randomIndex]
	}

	return string(otp)

}
