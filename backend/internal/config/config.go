package config

import (
	"os"
	"time"
)

// Config アプリケーション設定
type Config struct {
	JWTSecretKey        string
	JWTAccessTokenExpiry time.Duration
	JWTRefreshTokenExpiry time.Duration
	DBPath              string
	Port                string
}

// LoadConfig 環境変数から設定を読み込む
func LoadConfig() *Config {
	config := &Config{
		JWTSecretKey:        getEnv("JWT_SECRET_KEY", "your-super-secret-jwt-key-change-this-in-production"),
		JWTAccessTokenExpiry: getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", "1h"),
		JWTRefreshTokenExpiry: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", "168h"), // 7 days
		DBPath:              getEnv("DB_PATH", "./todo_app.db"),
		Port:                getEnv("PORT", "8080"),
	}

	return config
}

// getEnv 環境変数を取得し、デフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv 環境変数から時間を取得し、デフォルト値を返す
func getDurationEnv(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		// デフォルト値でパースを試す
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}
