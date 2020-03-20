package config

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

// GetEnv returns environment variable with fallback/default
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// CheckConfig ensures environment has been configured with GCP API access credentials
func CheckConfig() {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Error("*Required* key ID not set for Google application access.")
	}
	if os.Getenv("PROJECT") == "" {
		log.Error("*Required* project name not set for Google application access.")
	}
}
