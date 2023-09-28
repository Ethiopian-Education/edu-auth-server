package auth

import (
	"net/http"

	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/gin-gonic/gin"
)

func SignUpHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.Response{Message: "Successfully signed up", Success: true})
	}
}
