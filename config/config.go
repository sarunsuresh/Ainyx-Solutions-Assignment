package config

import (
	"fmt"
	"os"
	"time"
	 "strconv"
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
	EmailServiceAddr string
    EmailServicePort string
    EmailFailureRate float64
    CBMaxFailures    int
    CBResetTimeout   time.Duration
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
		EmailServiceAddr: getEnv("EMAIL_SERVICE_ADDR", "localhost:50051"),
        EmailServicePort: getEnv("EMAIL_SERVICE_PORT", "50051"),
        EmailFailureRate: getEnvFloat("EMAIL_FAILURE_RATE", 0.2),
        CBMaxFailures:    getEnvInt("CB_MAX_FAILURES", 5),
        CBResetTimeout:   getEnvDuration("CB_RESET_TIMEOUT", 30*time.Second),
		
	}
}

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



func getEnvFloat(key string, fallback float64) float64 {
    if v := os.Getenv(key); v != "" {
        if f, err := strconv.ParseFloat(v, 64); err == nil {
            return f
        }
    }
    return fallback
}

func getEnvInt(key string, fallback int) int {
    if v := os.Getenv(key); v != "" {
        if i, err := strconv.Atoi(v); err == nil {
            return i
        }
    }
    return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
    if v := os.Getenv(key); v != "" {
        if d, err := time.ParseDuration(v); err == nil {
            return d
        }
    }
    return fallback
}
