package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config contains the runtime settings required by the API process.
type Config struct {
	HTTPAddress string
	DatabaseURL string
}

// Load reads runtime settings from the process environment.
func Load() (Config, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL must be set")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return Config{}, fmt.Errorf("PORT must be a valid TCP port, got %q", port)
	}

	return Config{
		HTTPAddress: ":" + port,
		DatabaseURL: databaseURL,
	}, nil
}
