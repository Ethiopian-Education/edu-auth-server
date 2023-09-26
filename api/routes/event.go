package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addEventRoutes(rg *gin.RouterGroup) {
	event := rg.Group("/event")

	event.POST("/check_validity", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Requested item is active"})
	})
}