package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvVars loads all the env vars from the files depending on the
// environment. We use those for config files like the connection to the
// database
func LoadEnvVars() {
	appEnv := os.Getenv("APP_ENV")

	// Sets to development if no APP_ENV was passed
	if appEnv == "" {
		appEnv = "development"
	}

	os.Setenv("APP_ENV", appEnv)

	var err error

	switch appEnv {
	case "production":
		err = godotenv.Load()
	case "test":
		err = godotenv.Load(".env.test")
	default:
		err = godotenv.Load(".env.development")
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
