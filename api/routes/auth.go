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

}
