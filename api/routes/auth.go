package routes

import (
	auth_handler "github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/auth"
	"github.com/gin-gonic/gin"
)

func addAuthRoutes(rg *gin.RouterGroup) {
	auth_ := rg.Group("/auth")

	auth_.POST("/login", auth_handler.LoginHandler())
}

