package main

import (
	"fmt"

	"github.com/Ethiopian-Education/edu-auth-server.git/api/routes"
	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	_ "github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	// Initialize routes
	port := fmt.Sprintf(":%s", config.PORT)
	routes.InitRoutes(router)

	logrus.Infof("Server ...")
	router.Run(port)
}