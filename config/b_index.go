package config

import (
	"os"

	"github.com/Ethiopian-Education/edu-auth-server.git/crypto/parser"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	PORT        string
	TOKEN       string
	PRIVATE_KEY []byte
	PUBLIC_KEY  []byte
)

// Init function is invoked before main goroutine starts ... init functions are invoked based on the order they qued or based on their folder alphabetical orders
func init() {
	var err error
	// os.env
	godotenv.Load()

	logrus.Info("ENV loader")

	PORT = os.Getenv("PORT")

	PRIVATE_KEY, PUBLIC_KEY, err = parser.ReadKeys("./private.pem", "./public.pem")
	if err != nil {
		os.Exit(1)
	}

}
