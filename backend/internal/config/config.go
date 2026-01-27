package config

import (
	"os"
	"time"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name  string
	Env   string
	Debug bool
	Port  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// Load returns application configuration from environment
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:  getEnv("APP_NAME", "GoFiver"),
			Env:   getEnv("APP_ENV", "development"),
			Debug: getEnv("APP_DEBUG", "true") == "true",
			Port:  getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "mysql"),
			Port:     getEnv("DB_PORT", "3306"),
			Database: getEnv("DB_DATABASE", "gofiver"),
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", "secret"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "redis"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "secret"),
			Expiry: parseDuration(getEnv("JWT_EXPIRY", "24h")),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
