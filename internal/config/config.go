package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	ServerHost string
	ServerPort string

	JWTSecret     string
	SessionSecret string
	SessionMaxAge string

	RedisHost     string
	RedisPort     string
	RedisPassword string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		return nil
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "auth_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "3000"),

		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		SessionSecret: getEnv("SESSION_SECRET", "your-session-secret"),
		SessionMaxAge: getEnv("SESSION_MAX_AGE", "3600"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
