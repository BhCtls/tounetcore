package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment    string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiration  time.Duration
	NKeyExpiration time.Duration
	PushDeerAPI    string
	ServerPort     string
}

func LoadConfig() *Config {
	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		DatabaseURL:    getEnv("DATABASE_URL", "sqlite://./tounetcore.db"),
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpiration:  getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		NKeyExpiration: getDurationEnv("NKEY_EXPIRATION", 15*time.Minute),
		PushDeerAPI:    getEnv("PUSHDEER_API", "https://api2.pushdeer.com/message/push"),
		ServerPort:     getEnv("PORT", "44544"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
