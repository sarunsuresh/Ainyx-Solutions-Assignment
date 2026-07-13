package config

import (
	"fmt"
	"os"
)

// Config holds all environment-driven settings.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	JWTSecret  string
    JWTExpHours int
	RedisURL          string
    RateLimitRequests int    
    RateLimitWindow   int    
    WSAllowedOrigins  string
}

// Load reads config from environment variables, with sensible defaults.
func Load() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "userdb"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
		JWTSecret:   getEnv("JWT_SECRET", "changeme-secret"),
        JWTExpHours: 24,
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379"),
        RateLimitRequests: 10,  
        RateLimitWindow:   60,  
        WSAllowedOrigins:  getEnv("WS_ALLOWED_ORIGINS", "*"),
		
	}
}

// DSN builds the Postgres connection string.
func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
