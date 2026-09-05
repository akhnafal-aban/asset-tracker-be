package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	Port   string
	Env    string
	DBPath string
}

// Load reads configuration from the environment.
// It prioritizes standard library first, keeping it simple.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "asset_tracker.db"
	}

	return &Config{
		Port:   port,
		Env:    env,
		DBPath: dbPath,
	}
}
