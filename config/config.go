package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	AI       AIConfig
	CORS     CORSConfig
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type AIConfig struct {
	ServiceURL string
	Timeout    time.Duration
}

type CORSConfig struct {
	AllowedOrigins string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, continue without it
	}

	config := &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://username:password@localhost:5432/maths_solution_db?sslmode=disable"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		Server: ServerConfig{
			Port:    getEnv("PORT", "8000"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		AI: AIConfig{
			ServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:5000"),
			Timeout:    time.Duration(getEnvAsInt("AI_SERVICE_TIMEOUT", 30)) * time.Second,
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
