package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig
	Logging  LoggingConfig
	Database DatabaseConfig
}
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}
type LoggingConfig struct {
	Level string
}
type DatabaseConfig struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnvInt("APP_SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("APP_SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvDuration("APP_SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvDuration("APP_SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Logging: LoggingConfig{
			Level: getEnvString("APP_LOGGING_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Type:     getEnvString("APP_DATABASE_TYPE", "memory"),
			Host:     getEnvString("APP_DATABASE_HOST", "localhost"),
			Port:     getEnvInt("APP_DATABASE_PORT", 5432),
			User:     getEnvString("APP_DATABASE_USER", "postgres"),
			Password: getEnvString("APP_DATABASE_PASSWORD", "postgres"),
			DBName:   getEnvString("APP_DATABASE_DBNAME", "transaction_routine"),
			SSLMode:  getEnvString("APP_DATABASE_SSLMODE", "disable"),
		},
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
