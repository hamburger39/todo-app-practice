package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	JWTSecret      string
	JWTExpiryHours int
	DatabasePath   string
	Environment    string
}

func Load() *Config {
	config := &Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		DatabasePath:   getEnv("DATABASE_PATH", "./todo.db"),
		Environment:    getEnv("ENVIRONMENT", "development"),
	}

	// JWTシークレットが設定されていない場合は生成
	if config.JWTSecret == "" {
		config.JWTSecret = generateRandomSecret()
		if config.Environment == "production" {
			log.Println("WARNING: Using random JWT secret in production. Set JWT_SECRET environment variable.")
		}
	}

	return config
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

func generateRandomSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Failed to generate random secret:", err)
	}
	return hex.EncodeToString(bytes)
}
