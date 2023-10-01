package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	envs     []string
	messages []string
)

func init() {

	godotenv.Load()

	envs = []string{"PORT", "HASURA_GRAPHQL_ADMIN_SECRET"}
	// os.LookupEnv("key")
	for _, val := range envs {
		if el, exist := os.LookupEnv(val); !exist || el == "" {
			// env exist but its value is empty or it's not exist at all.
			messages = append(messages, val)
		}
	}

	if len(messages) > 0 {
		logrus.Errorf("Error encountered when checking those env vars ->>>>>> %v", strings.Join(messages, ", "))
		// panic("Error encountered when checking environmental variables")
		os.Exit(1)
	}
}
