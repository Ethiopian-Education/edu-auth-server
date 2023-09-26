package routes

import (
	"github.com/gin-gonic/gin"
	auth_handler "github.com/minilikmila/edu-auth-server.git/api/handlers/auth"
)

func addAuthRoutes(rg *gin.RouterGroup) {
	auth_ := rg.Group("/auth")

	auth_.POST("/login", auth_handler.LoginHandler())
}

