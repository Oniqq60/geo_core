package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Webhook  WebhookConfig
	Stats    StatsConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type WebhookConfig struct {
	URL           string
	RetryAttempts int
	RetryDelay    time.Duration
	Timeout       time.Duration
}

type StatsConfig struct {
	TimeWindowMinutes int
}

type AuthConfig struct {
	APIKey string
}

func Load() (*Config, error) {
	// Загружаем .env файл, если он существует
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "geo_system"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Webhook: WebhookConfig{
			URL:           getEnv("WEBHOOK_URL", "http://localhost:9090/webhook"),
			RetryAttempts: getEnvAsInt("WEBHOOK_RETRY_ATTEMPTS", 3),
			RetryDelay:    getEnvAsDuration("WEBHOOK_RETRY_DELAY", 5*time.Second),
			Timeout:       getEnvAsDuration("WEBHOOK_TIMEOUT", 10*time.Second),
		},
		Stats: StatsConfig{
			TimeWindowMinutes: getEnvAsInt("STATS_TIME_WINDOW_MINUTES", 60),
		},
		Auth: AuthConfig{
			APIKey: getEnv("API_KEY", "default-api-key-change-in-production"),
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
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
