package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_FromEnv(t *testing.T) {
	// Arrange
	expectedKey := "test-api-key-from-env"
	os.Setenv("QWEATHER_API_KEY", expectedKey)
	defer os.Unsetenv("QWEATHER_API_KEY")

	// Act
	config, err := Load()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedKey, config.QWeather.APIKey)
	assert.Equal(t, "devapi.qweather.com", config.QWeather.APIHost)
}

func TestLoad_CustomAPIHost(t *testing.T) {
	// Arrange
	expectedKey := "test-api-key"
	expectedHost := "api.qweather.com"
	os.Setenv("QWEATHER_API_KEY", expectedKey)
	os.Setenv("QWEATHER_API_HOST", expectedHost)
	defer os.Unsetenv("QWEATHER_API_KEY")
	defer os.Unsetenv("QWEATHER_API_HOST")

	// Act
	config, err := Load()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedKey, config.QWeather.APIKey)
	assert.Equal(t, expectedHost, config.QWeather.APIHost)
}

func TestLoad_FromFile(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "awesome-skills", "qweather")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	expectedKey := "test-api-key-from-file"
	apiKeyFile := filepath.Join(configDir, "api_key")
	require.NoError(t, os.WriteFile(apiKeyFile, []byte(expectedKey+"\n"), 0600))

	// Override HOME for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Ensure no env var interferes
	os.Unsetenv("QWEATHER_API_KEY")

	// Act
	config, err := Load()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedKey, config.QWeather.APIKey)
}

func TestLoad_FileWithWhitespace(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "awesome-skills", "qweather")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	expectedKey := "test-api-key-trimmed"
	apiKeyFile := filepath.Join(configDir, "api_key")
	// Add whitespace and newlines
	require.NoError(t, os.WriteFile(apiKeyFile, []byte("  "+expectedKey+"  \n\n"), 0600))

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	os.Unsetenv("QWEATHER_API_KEY")

	// Act
	config, err := Load()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedKey, config.QWeather.APIKey)
}

func TestLoad_EnvPriority(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "awesome-skills", "qweather")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	fileKey := "file-key"
	envKey := "env-key"

	apiKeyFile := filepath.Join(configDir, "api_key")
	require.NoError(t, os.WriteFile(apiKeyFile, []byte(fileKey), 0600))

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	os.Setenv("QWEATHER_API_KEY", envKey)
	defer os.Unsetenv("QWEATHER_API_KEY")

	// Act
	config, err := Load()

	// Assert - environment variable should take priority
	require.NoError(t, err)
	assert.Equal(t, envKey, config.QWeather.APIKey)
}

func TestLoad_NoAPIKey(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	os.Unsetenv("QWEATHER_API_KEY")

	// Act
	config, err := Load()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "QWeather API key not found")
	assert.Nil(t, config)
}

func TestLoad_EmptyAPIKeyFile(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "awesome-skills", "qweather")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	apiKeyFile := filepath.Join(configDir, "api_key")
	require.NoError(t, os.WriteFile(apiKeyFile, []byte("   \n  \n"), 0600))

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	os.Unsetenv("QWEATHER_API_KEY")

	// Act
	config, err := Load()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API key file is empty")
	assert.Nil(t, config)
}

func TestLoad_XDGConfigHome(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	xdgConfig := filepath.Join(tmpDir, "custom-config")
	configDir := filepath.Join(xdgConfig, "awesome-skills", "qweather")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	expectedKey := "test-xdg-key"
	apiKeyFile := filepath.Join(configDir, "api_key")
	require.NoError(t, os.WriteFile(apiKeyFile, []byte(expectedKey), 0600))

	os.Setenv("XDG_CONFIG_HOME", xdgConfig)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	os.Unsetenv("QWEATHER_API_KEY")

	// Act
	config, err := Load()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedKey, config.QWeather.APIKey)
}

func TestEnsureConfigDir(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Act
	err := EnsureConfigDir()

	// Assert
	require.NoError(t, err)

	// Verify directory exists
	expectedDir := filepath.Join(tmpDir, ".config", "awesome-skills", "qweather")
	info, err := os.Stat(expectedDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Verify permissions
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
}

func TestEnsureConfigDir_AlreadyExists(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "awesome-skills", "qweather")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Act
	err := EnsureConfigDir()

	// Assert - should not error if directory already exists
	require.NoError(t, err)
}

func TestEnsureConfigDir_XDGConfigHome(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	xdgConfig := filepath.Join(tmpDir, "custom-config")

	os.Setenv("XDG_CONFIG_HOME", xdgConfig)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	// Act
	err := EnsureConfigDir()

	// Assert
	require.NoError(t, err)

	// Verify directory exists in XDG location
	expectedDir := filepath.Join(xdgConfig, "awesome-skills", "qweather")
	info, err := os.Stat(expectedDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}
