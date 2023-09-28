package config

import (
	"os"

	"github.com/Ethiopian-Education/edu-auth-server.git/crypto/parser"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	PORT                   string
	TOKEN                  string
	PRIVATE_KEY            []byte
	PUBLIC_KEY             []byte
	TWILIO_ACCOUNT_SID     string
	TWILIO_AUTH_TOKEN      string
	TWILIO_PHONE_NUMBER    string
	TWILIO_WHATSAPP_NUMBER string
	SMTP_PORT              string
	SMTP_HOST              string
	SMTP_USERNAME          string
	SMTP_PASSWORD          string
)

// Init function is invoked before main goroutine starts ... init functions are invoked based on the order they qued or based on their folder alphabetical orders
func init() {
	var err error
	// os.env
	godotenv.Load()

	logrus.Info("ENV loader")

	PORT = os.Getenv("PORT")

	TWILIO_ACCOUNT_SID = os.Getenv("TWILIO_ACCOUNT_SID")

	TWILIO_AUTH_TOKEN = os.Getenv("TWILIO_AUTH_TOKEN")

	TWILIO_PHONE_NUMBER = os.Getenv("TWILIO_PHONE_NUMBER")

	TWILIO_WHATSAPP_NUMBER = os.Getenv("TWILIO_WHATSAPP_NUMBER")

	SMTP_PORT = os.Getenv("SMTP_PORT")

	SMTP_HOST = os.Getenv("SMTP_HOST")

	SMTP_USERNAME = os.Getenv("SMTP_USERNAME")

	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")

	PRIVATE_KEY, PUBLIC_KEY, err = parser.ReadKeys("./private.pem", "./public.pem")
	if err != nil {
		os.Exit(1)
	}

}
