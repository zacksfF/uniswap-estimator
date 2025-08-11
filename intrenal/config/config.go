package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host string
	Port int

	// Blockchain Config
	EthereumRPCURL string

	// Performance settings
	RequestTimeout time.Duration
	MaxConnections int

	// Environment
	Environment string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Server settings
	config.Host = getEnvOrDefault("HOST", "localhost")

	port, err := strconv.Atoi(getEnvOrDefault("PORT", "1337"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %v", err)
	}
	config.Port = port

	// Blockchain settings - REQUIRED for 1inch assignment
	config.EthereumRPCURL = os.Getenv("ETHEREUM_RPC_URL")
	if config.EthereumRPCURL == "" {
		return nil, fmt.Errorf("ETHEREUM_RPC_URL is required")
	}

	// Performance settings
	timeout, _ := strconv.Atoi(getEnvOrDefault("REQUEST_TIMEOUT", "10"))
	config.RequestTimeout = time.Duration(timeout) * time.Second

	maxConn, _ := strconv.Atoi(getEnvOrDefault("MAX_CONNECTIONS", "100"))
	config.MaxConnections = maxConn

	// Environment
	config.Environment = getEnvOrDefault("ENV", "development")

	return config, nil
}

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsProduction checks if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}