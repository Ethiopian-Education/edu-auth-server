package routes

import (
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitRoutes(router *gin.Engine) {
	phone := "00918492083"
	phone, isValid := utils.ValidatePhone(phone)
	if !isValid {
		logrus.Error("Invalid phone number ")
	}
	logrus.Info("PHONE NEWLY generated : ", phone)
	// gin.SetMode("release")
  logrus.Infoln("Inside Router Package and initialize the route handlers...")
  getRoutes(router)
	 
}

func getRoutes (router *gin.Engine) {
	v1 := router.Group("/v1")
	addAuthRoutes(v1)
	addEventRoutes(v1)
}