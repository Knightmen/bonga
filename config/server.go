package config

import (
	"os"
)

// GetSwaggerHost returns the host for Swagger documentation
func GetSwaggerHost() string {
	host := os.Getenv("SWAGGER_HOST")
	if host == "" {
		host = "localhost:8080"
	}
	return host
} 