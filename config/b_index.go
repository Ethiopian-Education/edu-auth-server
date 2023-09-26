package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (

	PORT string
	TOKEN string
)
// Init function is invoked before main goroutine starts ... init functions are invoked based on the order they qued or based on their folder alphabetical orders
func init() {
	// os.env
	godotenv.Load()

	logrus.Info("ENV loader")
	PORT = os.Getenv("PORT")
}