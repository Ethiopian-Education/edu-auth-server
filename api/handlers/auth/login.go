package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minilikmila/edu-auth-server.git/model"
)

func LoginHandler() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.Response{
			Message: "Successful login <here-token>",
			Success: true,
		})
	}
}