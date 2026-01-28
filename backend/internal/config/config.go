package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Cache    CacheConfig
}

type AppConfig struct {
	Name  string
	Env   string
	Debug bool
	Port  string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	Database        string
	Username        string
	Password        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	Host        string
	Port        string
	Password    string
	PoolSize    int
	MinIdleConn int
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type CacheConfig struct {
	Enabled       bool
	ListTTL       time.Duration
	CountTTL      time.Duration
	DetailTTL     time.Duration
	DefaultPerPage int
	MaxPerPage    int
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:  getEnv("APP_NAME", "GoFiver"),
			Env:   getEnv("APP_ENV", "development"),
			Debug: getEnv("APP_DEBUG", "true") == "true",
			Port:  getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "mysql"),
			Port:            getEnv("DB_PORT", "3306"),
			Database:        getEnv("DB_DATABASE", "gofiver"),
			Username:        getEnv("DB_USERNAME", "root"),
			Password:        getEnv("DB_PASSWORD", "secret"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: parseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m")),
			ConnMaxIdleTime: parseDuration(getEnv("DB_CONN_MAX_IDLE_TIME", "1m")),
		},
		Redis: RedisConfig{
			Host:        getEnv("REDIS_HOST", "redis"),
			Port:        getEnv("REDIS_PORT", "6379"),
			Password:    getEnv("REDIS_PASSWORD", ""),
			PoolSize:    getEnvInt("REDIS_POOL_SIZE", 100),
			MinIdleConn: getEnvInt("REDIS_MIN_IDLE_CONN", 10),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "secret"),
			Expiry: parseDuration(getEnv("JWT_EXPIRY", "24h")),
		},
		Cache: CacheConfig{
			Enabled:       getEnv("CACHE_ENABLED", "true") == "true",
			ListTTL:       parseDuration(getEnv("CACHE_LIST_TTL", "30s")),
			CountTTL:      parseDuration(getEnv("CACHE_COUNT_TTL", "60s")),
			DetailTTL:     parseDuration(getEnv("CACHE_DETAIL_TTL", "5m")),
			DefaultPerPage: getEnvInt("CACHE_DEFAULT_PER_PAGE", 20),
			MaxPerPage:    getEnvInt("CACHE_MAX_PER_PAGE", 100),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
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
