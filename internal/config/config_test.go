package config_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Melsoft-Games/ant-watcher/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestLoadConfigFile checks the correctness of loading the configuration from a temporary file
func TestLoadConfigFile(t *testing.T) {
	// Create a temporary file with test configuration
	configData := `{
		"github_token": "test_file_token",
		"github_api_url": "https://file-test.github.com",
		"webhook_address": "3.2.1.4",
		"webhook_port": "24132",
		"webhook_secret": "test_file_secret",
		"admin_address": "127.0.0.6",
		"admin_port": "321",
		"log_level": "WARNING",
		"push_metrics_url": "http://victoriametrics-file:8430",
		"memory_ttl": "12h",
		"fetch_history": "10h15m",
		"memory_limit": "100g",
		"disable_admin_server": "true",
		"disable_api": "false",
		"metrics_address": "5.6.7.8",
		"metrics_port": "9100"
	}`
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Удалить файл после теста

	// Write data to the temporary file
	if _, err := tmpFile.Write([]byte(configData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the file to complete the write
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Override the configuration path
	os.Setenv("CONFIG_FILE_PATH", tmpFile.Name())

	// Load the configuration
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "test_file_token", cfg.GitHubToken)
		assert.Equal(t, "https://file-test.github.com", cfg.GitHubAPIURL)
		assert.Equal(t, "3.2.1.4", cfg.WebhookAddress)
		assert.Equal(t, "24132", cfg.WebhookPort)
		assert.Equal(t, "test_file_secret", cfg.WebhookSecret)
		assert.Equal(t, "127.0.0.6", cfg.AdminAddress)
		assert.Equal(t, "321", cfg.AdminPort)
		assert.Equal(t, "WARNING", cfg.LogLevel)
		assert.Equal(t, "http://victoriametrics-file:8430", cfg.PushMetricsUrl)
		assert.Equal(t, "12h", cfg.MemoryTTL)
		assert.Equal(t,
			func() time.Duration { d, _ := time.ParseDuration("12h"); return d }(),
			cfg.MemoryTTLTime)
		assert.Equal(t, "10h15m", cfg.FetchHistory)
		assert.Equal(t,
			func() time.Duration { d, _ := time.ParseDuration("10h15m"); return d }(),
			cfg.FetchHistoryTime)
		assert.Equal(t, "100g", cfg.MemoryLimit)
		assert.Equal(t, uint64(100*1024*1024*1024), cfg.MemoryLimitBytes)
		assert.Equal(t, "true", cfg.DisableAdminServer)
		assert.Equal(t, "false", cfg.DisableAPI)
		assert.Equal(t, "5.6.7.8", cfg.MetricsAddress)
		assert.Equal(t, "9100", cfg.MetricsPort)
	} else {
		t.Errorf("Config from test-tmp file is nil")
	}
}

// TestLoadDefaultConfig checks the correctness of loading the default configuration
func TestLoadDefaultConfig(t *testing.T) {
	// LoadConfig
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "", cfg.GitHubToken)
	assert.Equal(t, "https://api.github.com", cfg.GitHubAPIURL)
	assert.Equal(t, "0.0.0.0", cfg.WebhookAddress)
	assert.Equal(t, "8080", cfg.WebhookPort)
	assert.Equal(t, "", cfg.WebhookSecret)
	assert.Equal(t, "127.0.0.1", cfg.AdminAddress)
	assert.Equal(t, "8081", cfg.AdminPort)
	assert.Equal(t, "INFO", cfg.LogLevel)
	assert.Equal(t, "", cfg.PushMetricsUrl)
	assert.Equal(t, "15m", cfg.MemoryTTL)
	assert.Equal(t,
		func() time.Duration { d, _ := time.ParseDuration("15m"); return d }(),
		cfg.MemoryTTLTime)
	assert.Equal(t, "15m", cfg.FetchHistory)
	assert.Equal(t,
		func() time.Duration { d, _ := time.ParseDuration("15m"); return d }(),
		cfg.FetchHistoryTime)
	assert.Equal(t, "0", cfg.MemoryLimit)
	assert.Equal(t,
		func() uint64 { u, _ := strconv.ParseUint("0", 10, 64); return u }(),
		cfg.MemoryLimitBytes)
	assert.Equal(t, "0.0.0.0", cfg.MetricsAddress)
	assert.Equal(t, "3000", cfg.MetricsPort)
}

func TestLoadEnvConfig(t *testing.T) {
	// Set environment variables
	os.Setenv("ADMIN_ADDRESS", "127.0.0.4")
	os.Setenv("ADMIN_PORT", "1234")
	os.Setenv("GITHUB_API_URL", "https://ent.github.com")
	os.Setenv("GITHUB_TOKEN", "ENV_TEST_TOKEN")
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("PUSH_METRICS_URL", "http://victoriametrics-env:8429")
	os.Setenv("WEBHOOK_ADDRESS", "1.2.3.4")
	os.Setenv("WEBHOOK_PORT", "4321")
	os.Setenv("WEBHOOK_SECRET", "ENV_TEST_SECRET")
	os.Setenv("MEMORY_TTL", "100h")
	os.Setenv("FETCH_HISTORY", "101h")
	os.Setenv("MEMORY_LIMIT", "10Gb")
	os.Setenv("DISABLE_ADMIN_SERVER", "true")
	os.Setenv("DISABLE_API", "true")
	os.Setenv("METRICS_ADDRESS", "9.8.7.6")
	os.Setenv("METRICS_PORT", "9200")

	// Load the configuration
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "ENV_TEST_TOKEN", cfg.GitHubToken)
	assert.Equal(t, "https://ent.github.com", cfg.GitHubAPIURL)
	assert.Equal(t, "127.0.0.4", cfg.AdminAddress)
	assert.Equal(t, "1234", cfg.AdminPort)
	assert.Equal(t, "https://ent.github.com", cfg.GitHubAPIURL)
	assert.Equal(t, "ENV_TEST_TOKEN", cfg.GitHubToken)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, "http://victoriametrics-env:8429", cfg.PushMetricsUrl)
	assert.Equal(t, "1.2.3.4", cfg.WebhookAddress)
	assert.Equal(t, "4321", cfg.WebhookPort)
	assert.Equal(t, "ENV_TEST_SECRET", cfg.WebhookSecret)
	assert.Equal(t, "100h", cfg.MemoryTTL)
	assert.Equal(t,
		func() time.Duration { d, _ := time.ParseDuration("100h"); return d }(),
		cfg.MemoryTTLTime)
	assert.Equal(t, "101h", cfg.FetchHistory)
	assert.Equal(t,
		func() time.Duration { d, _ := time.ParseDuration("101h"); return d }(),
		cfg.FetchHistoryTime)
	assert.Equal(t, "10Gb", cfg.MemoryLimit)
	assert.Equal(t, uint64(10*1024*1024*1024), cfg.MemoryLimitBytes)
	assert.Equal(t, "true", cfg.DisableAdminServer)
	assert.Equal(t, "true", cfg.DisableAPI)
	assert.Equal(t, "9.8.7.6", cfg.MetricsAddress)
	assert.Equal(t, "9200", cfg.MetricsPort)
}
