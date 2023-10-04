package routes

import (
	auth_handler "github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/auth"
	"github.com/gin-gonic/gin"
)

func addAuthRoutes(rg *gin.RouterGroup) {
	auth_ := rg.Group("/auth")

	auth_.POST("/login", auth_handler.LoginHandler())
	auth_.POST("/signup", auth_handler.SignUpHandler())
	auth_.POST("/verify_phone", auth_handler.VerifyPhoneNumber())
	auth_.POST("/verify_2fa_auth", auth_handler.Authenticate2FA())
	auth_.POST("/forgot_password", auth_handler.ForgotPassword())
	auth_.POST("/reset_password", auth_handler.ResetPassword())
	auth_.POST("/update_password", auth_handler.UpdatePassword())
	auth_.POST("/resend_opt", auth_handler.ResendOTP())

}
