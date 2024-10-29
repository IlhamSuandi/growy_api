package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/joho/godotenv"
)

var (
	APP_ENV                    string
	DB_HOST                    string
	DB_USER                    string
	DB_PASSWORD                string
	DB_NAME                    string
	DB_PORT                    string
	DB_SSLMODE                 string
	DB_TIMEZONE                string
	JWT_SECRET                 string
	ALLOWED_ORIGIN             string
	JWT_EXPIRATION             int
	JWT_REFRESH_EXPIRATION     int
	BASE_URL                   string
	URL                        string
	OAUTH2_CLIENT_ID           string
	OAUTH2_CLIENT_SECRET       string
	OAUTH2_ClIENT_REDIRECT_URL string
)

func init() {
	loadEnv()

	jwtExpiration, err := strconv.Atoi(GetEnv("JWT_EXPIRATION", "604800"))
	if err != nil {
		jwtExpiration = 604800
	}

	jwtRefreshExpiration, err := strconv.Atoi(GetEnv("JWT_REFRESH_EXPIRATION", "2592000"))
	if err != nil {
		jwtRefreshExpiration = 2592000
	}

	APP_ENV = GetEnv("APP_ENV", "development")
	DB_HOST = GetEnv("DB_HOST", "localhost")
	DB_USER = GetEnv("DB_USER", "root")
	DB_PASSWORD = GetEnv("DB_PASSWORD", "")
	DB_NAME = GetEnv("DB_NAME", "business_assistant")
	DB_PORT = GetEnv("DB_PORT", "5432")
	DB_SSLMODE = GetEnv("DB_SSLMODE", "disable")
	DB_TIMEZONE = GetEnv("DB_TIMEZONE", "Asia/Jakarta")
	JWT_SECRET = GetEnv("JWT_SECRET", "secret")
	ALLOWED_ORIGIN = GetEnv("ALLOWED_ORIGIN", "")
	JWT_EXPIRATION = jwtExpiration
	JWT_REFRESH_EXPIRATION = jwtRefreshExpiration
	BASE_URL = GetEnv("BASE_URL", "http://localhost:5000")
	URL = GetEnv("URL", "http://localhost:5000/api/v1")
	OAUTH2_CLIENT_ID = GetEnv("OAUTH2_CLIENT_ID", "")
	OAUTH2_CLIENT_SECRET = GetEnv("OAUTH2_CLIENT_SECRET", "")
	OAUTH2_ClIENT_REDIRECT_URL = GetEnv("OAUTH2_REDIRECT_URL", "")
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func loadEnv() {
	APP_ENV = os.Getenv("APP_ENV")
	var env string

	switch APP_ENV {
	case "development":
		env = ".env.development"
	case "production":
		env = ".env.production"
	default:
		env = ".env.development"
	}

	// Define possible paths for the .env file
	configPaths := []string{
		fmt.Sprintf("./%s", env), // For app in current directory
		"../../../.env.test",     // For test folder
	}

	// Iterate over the possible paths
	for _, path := range configPaths {
		// Attempt to load .env file
		if err := godotenv.Load(path); err == nil {
			utils.Log.Infof("loaded env from %s", path)
			return
		}
	}

	utils.Log.Infof("failed to load env from any of: %v", configPaths)
}
