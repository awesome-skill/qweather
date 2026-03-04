package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	QWeather QWeatherConfig
}

// QWeatherConfig contains QWeather API configuration
type QWeatherConfig struct {
	APIKey  string
	APIHost string
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	// Load .env file if it exists
	// This will not error if .env doesn't exist, which is fine
	_ = godotenv.Load()

	cfg := &Config{
		QWeather: QWeatherConfig{
			APIHost: getEnvWithDefault("QWEATHER_API_HOST", "devapi.qweather.com"),
		},
	}

	// Try to load API key from environment variable first
	apiKey := os.Getenv("QWEATHER_API_KEY")
	if apiKey != "" {
		cfg.QWeather.APIKey = apiKey
		return cfg, nil
	}

	// Try to load from config file
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("get config directory: %w", err)
	}

	apiKeyPath := filepath.Join(configDir, "qweather", "api_key")
	apiKeyBytes, err := os.ReadFile(apiKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("QWeather API key not found. Please set QWEATHER_API_KEY environment variable or create %s", apiKeyPath)
		}
		return nil, fmt.Errorf("read API key file: %w", err)
	}

	cfg.QWeather.APIKey = strings.TrimSpace(string(apiKeyBytes))
	if cfg.QWeather.APIKey == "" {
		return nil, fmt.Errorf("API key file is empty: %s", apiKeyPath)
	}

	if os.Getenv("QWEATHER_API_HOST") == "" {
		apiHostPath := filepath.Join(configDir, "qweather", "api_host")
		if apiHostBytes, err := os.ReadFile(apiHostPath); err == nil {
			if host := strings.TrimSpace(string(apiHostBytes)); host != "" {
				cfg.QWeather.APIHost = host
			}
		}
	}

	return cfg, nil
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	// Check for XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "awesome-skill"), nil
	}

	// Fall back to ~/.config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", "awesome-skill"), nil
}

// getEnvWithDefault gets environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	qweatherDir := filepath.Join(configDir, "qweather")
	if err := os.MkdirAll(qweatherDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	return nil
}

// GetConfigDir returns the config directory path
func GetConfigDir() (string, error) {
	return getConfigDir()
}

// SetAPIKey saves the API key to the config file
func SetAPIKey(apiKey string) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	apiKeyPath := filepath.Join(configDir, "qweather", "api_key")
	if err := os.WriteFile(apiKeyPath, []byte(apiKey), 0600); err != nil {
		return fmt.Errorf("write API key file: %w", err)
	}

	return nil
}

// SetAPIHost saves the API host to the config file
func SetAPIHost(apiHost string) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	apiHostPath := filepath.Join(configDir, "qweather", "api_host")
	if err := os.WriteFile(apiHostPath, []byte(apiHost), 0600); err != nil {
		return fmt.Errorf("write API host file: %w", err)
	}

	return nil
}
