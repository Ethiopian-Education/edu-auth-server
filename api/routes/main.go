package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitRoutes(router *gin.Engine) {
	// gin.SetMode("release")
  logrus.Infoln("Inside Router Package and initialize the route handlers...")
  getRoutes(router)
	 
}

func getRoutes (router *gin.Engine) {
	v1 := router.Group("/v1")
	addAuthRoutes(v1)
	addEventRoutes(v1)
}