package auth

import "github.com/gin-gonic/gin"

type resendOTPBody struct {
	Email string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Type string `json:"type"`
}

func ResendOTP() gin.HandlerFunc{
	return func(ctx *gin.Context) {}
}