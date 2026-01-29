package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// Try loading .env first (for local development)
	envPath := filepath.Join(dir, "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		slog.Warn("failed to load .env file, trying .env.docker.local", "error", err, "path", envPath)
		// Try loading .env.docker.local as fallback
		dockerEnvPath := filepath.Join(dir, "..", ".env.docker.local")
		if err := godotenv.Load(dockerEnvPath); err != nil {
			// If both fail, it's okay - environment variables might already be set (e.g., by docker compose)
			slog.Info("env files not found, using environment variables from system", "tried", []string{envPath, dockerEnvPath})
		} else {
			envPath = dockerEnvPath
		}
	}

	// Build DATABASE_URL from individual DB variables if DATABASE_URL is not set
	if os.Getenv("DATABASE_URL") == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")
		dbUsername := os.Getenv("DB_USERNAME")
		dbPassword := os.Getenv("DB_PASSWORD")

		if dbHost != "" && dbPort != "" && dbName != "" && dbUsername != "" && dbPassword != "" {
			dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
				dbUsername, dbPassword, dbHost, dbPort, dbName)
			os.Setenv("DATABASE_URL", dbUrl)
			slog.Info("constructed DATABASE_URL from individual DB variables")
		}
	}

	// Map SERVER_PORT to PORT if PORT is not set
	if os.Getenv("PORT") == "" {
		if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
			os.Setenv("PORT", ":"+serverPort)
			slog.Info("mapped SERVER_PORT to PORT", "port", serverPort)
		}
	}

	slog.Info("config loaded successfully", "path", envPath)
	return nil
}
