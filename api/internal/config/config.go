package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddress     string
	DatabaseURL     string
	APIToken        string
	OpenAPISpecPath string
}

// Load reads and validates runtime settings from the process environment.
func Load() (Config, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL must be set")
	}

	apiToken := strings.TrimSpace(os.Getenv("API_TOKEN"))
	if apiToken == "" {
		return Config{}, fmt.Errorf("API_TOKEN must be set")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return Config{}, fmt.Errorf("PORT must be a valid TCP port, got %q", port)
	}
	openAPISpecPath := strings.TrimSpace(os.Getenv("OPENAPI_SPEC_PATH"))
	if openAPISpecPath == "" {
		openAPISpecPath = "../docs/api.yaml"
	}

	return Config{
		HTTPAddress:     ":" + port,
		DatabaseURL:     databaseURL,
		APIToken:        apiToken,
		OpenAPISpecPath: openAPISpecPath,
	}, nil
}
