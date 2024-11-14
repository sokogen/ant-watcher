package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains the processed configuration
type Config struct {
	AdminAddress       string        `json:"admin_address"`        // Address to listen for admin requests
	AdminPort          string        `json:"admin_port"`           // Port to listen for admin requests
	DisableAdminServer string        `json:"disable_admin_server"` // Turn off the admin server, only while starting the app
	MetricsAddress     string        `json:"metrics_address"`      // Address to listen for metrics requests
	MetricsPort        string        `json:"metrics_port"`         // Port to listen for metrics requests
	PushMetricsUrl     string        `json:"push_metrics_url"`     // Address to push metrics to Prometheus/VictoriaMetrics
	DisableAPI         string        `json:"disable_api"`          // Turn off the API server, only while starting the app
	WebhookAddress     string        `json:"webhook_address"`      // Address to listen for incoming webhooks
	WebhookPort        string        `json:"webhook_port"`         // Port to listen for incoming webhooks
	WebhookSecret      string        `json:"webhook_secret"`       // Secret key for webhook validation
	GitHubToken        string        `json:"github_token"`         // Token for GitHub API, if empty, the API will not be used
	GitHubAPIURL       string        `json:"github_api_url"`       // URL for GitHub API
	LogLevel           string        `json:"log_level"`            // Log level [DEBUG, INFO, WARN, ERROR, FATAL]
	MemoryTTL          string        `json:"memory_ttl"`           // Memory TTL in human-readable format
	MemoryTTLTime      time.Duration `json:"-"`                    // Time to live for objects in memory (computed, not from JSON)
	FetchHistory       string        `json:"fetch_history"`        // Time of previous events to fetch
	FetchHistoryTime   time.Duration `json:"-"`                    // Time of previous events to fetch (computed, not from JSON)
	MemoryLimit        string        `json:"memory_limit"`         // Memory limit in human-readable format
	MemoryLimitBytes   uint64        `json:"-"`                    // Memory limit in bytes (computed, not from JSON)
}

// for administration and tests purposes
type ConfigInterface interface {
	ReloadConfig() error
}

var configlog *log.Logger

const (
	defAdminAddress         = "127.0.0.1"
	defAdminPort            = "8081"
	defDisableAdminServer   = "false"
	defDisableAPI           = "false"
	defFetchHistory         = "15m"
	defGitHubAPIURL         = "https://api.github.com"
	defGitHubAppID          = ""
	defGitHubInstallationID = ""
	defLogLevel             = "INFO"
	defMemoryLimit          = "0"
	defMemoryTTL            = "15m"
	defMetricsAddress       = "0.0.0.0"
	defMetricsPort          = "3000"
	defPushMetricsUrl       = ""
	defWebhookAddress       = "0.0.0.0"
	defWebhookPort          = "8080"
	localLogBanner          = "CONFIG"
)

// local realization of the logger,
// need for solving problem of chicken-egg
func printConfigEvent(v ...interface{})                 { configlog.Println(v...) }
func printConfigEventf(format string, v ...interface{}) { configlog.Printf(format, v...) }

// getEnv fetches a string value from the environment or returns the default.
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// loadAndProcessConfig reads configuration from file, processes it, applies default values, and overrides with environment variables
func loadAndProcessConfig() (*Config, error) {
	// Default values, placed here because some of them are should be used as result of function
	// and we can't use them as default values in struct fields
	rawCfg := Config{
		AdminAddress:       defAdminAddress,
		MetricsAddress:     defMetricsAddress,
		MetricsPort:        defMetricsPort,
		AdminPort:          defAdminPort,
		WebhookAddress:     defWebhookAddress,
		WebhookPort:        defWebhookPort,
		GitHubAPIURL:       defGitHubAPIURL,
		LogLevel:           defLogLevel,
		MemoryTTL:          defMemoryTTL,
		MemoryLimit:        defMemoryLimit,
		PushMetricsUrl:     defPushMetricsUrl,
		FetchHistory:       defFetchHistory,
		DisableAdminServer: defDisableAdminServer,
		DisableAPI:         defDisableAPI,
	}

	configFilePath := getEnv("CONFIG_FILE_PATH", "config/config.json")
	cwd, _ := os.Getwd()
	printConfigEventf("Current working directory: %s", cwd)
	printConfigEventf("Attempting to read config file at: %s", configFilePath)

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		printConfigEventf("config file not found at %s: %v", configFilePath, err)
	} else {
		printConfigEventf("Config file found at %s, parsing...", configFilePath)
		if err := json.Unmarshal(data, &rawCfg); err != nil {
			return nil, fmt.Errorf("could not unmarshal config: %v", err)
		}
		printConfigEvent("Config file parsed successfully")
	}

	// Apply values from environment variables
	envVars := map[string]*string{
		"ADMIN_ADDRESS":        &rawCfg.AdminAddress,
		"ADMIN_PORT":           &rawCfg.AdminPort,
		"METRICS_ADDRESS":      &rawCfg.MetricsAddress,
		"METRICS_PORT":         &rawCfg.MetricsPort,
		"GITHUB_API_URL":       &rawCfg.GitHubAPIURL,
		"GITHUB_TOKEN":         &rawCfg.GitHubToken,
		"LOG_LEVEL":            &rawCfg.LogLevel,
		"PUSH_METRICS_URL":     &rawCfg.PushMetricsUrl,
		"WEBHOOK_ADDRESS":      &rawCfg.WebhookAddress,
		"WEBHOOK_PORT":         &rawCfg.WebhookPort,
		"WEBHOOK_SECRET":       &rawCfg.WebhookSecret,
		"DISABLE_ADMIN_SERVER": &rawCfg.DisableAdminServer,
		"DISABLE_API":          &rawCfg.DisableAPI,
		"MEMORY_TTL":           &rawCfg.MemoryTTL,
		"FETCH_HISTORY":        &rawCfg.FetchHistory,
		"MEMORY_LIMIT":         &rawCfg.MemoryLimit,
	}

	for key, ptr := range envVars {
		*ptr = getEnv(key, *ptr)
	}

	// Validate the configuration
	if net.ParseIP(rawCfg.AdminAddress) == nil {
		return nil, fmt.Errorf("invalid AdminAddress: %s", rawCfg.AdminAddress)
	}

	if _, err := strconv.Atoi(rawCfg.AdminPort); err != nil {
		return nil, fmt.Errorf("invalid AdminPort: %s", rawCfg.AdminPort)
	}

	if net.ParseIP(rawCfg.MetricsAddress) == nil {
		return nil, fmt.Errorf("invalid MetricsAddress: %s", rawCfg.MetricsAddress)
	}

	if _, err := strconv.Atoi(rawCfg.MetricsPort); err != nil {
		return nil, fmt.Errorf("invalid MetricsPort: %s", rawCfg.MetricsPort)
	}

	if _, err := url.Parse(rawCfg.GitHubAPIURL); err != nil {
		return nil, fmt.Errorf("invalid GitHubAPIURL: %s", rawCfg.GitHubAPIURL)
	}

	if rawCfg.PushMetricsUrl != "" {
		if _, err := url.ParseRequestURI(rawCfg.PushMetricsUrl); err != nil {
			return nil, fmt.Errorf("invalid PushMetricsUrl: %s", rawCfg.PushMetricsUrl)
		}
	}

	if net.ParseIP(rawCfg.WebhookAddress) == nil {
		return nil, fmt.Errorf("invalid WebhookAddress: %s", rawCfg.WebhookAddress)
	}

	if _, err := strconv.Atoi(rawCfg.WebhookPort); err != nil {
		return nil, fmt.Errorf("invalid WebhookPort: %s", rawCfg.WebhookPort)
	}

	// Processing MemoryTTL, FetchHistory and MemoryLimit values
	rawCfg.MemoryTTLTime, err = time.ParseDuration(rawCfg.MemoryTTL)
	if err != nil {
		return nil, fmt.Errorf("invalid MemoryTTL: %v", err)
	}

	rawCfg.FetchHistoryTime, err = time.ParseDuration(rawCfg.FetchHistory)
	if err != nil {
		return nil, fmt.Errorf("invalid FetchHistory: %v", err)
	}

	rawCfg.MemoryLimitBytes, err = parseSize(rawCfg.MemoryLimit)
	if err != nil {
		return nil, fmt.Errorf("invalid MemoryLimit: %v", err)
	}

	return &rawCfg, nil
}

// LoadConfig loads the configuration initially
func LoadConfig() (*Config, error) {
	configlog = log.New(os.Stdout, localLogBanner+": ", log.LstdFlags)
	cfg, err := loadAndProcessConfig()
	if err != nil {
		return nil, err
	}
	printConfigEvent("Configuration loaded successfully")
	return cfg, nil
}

// ReloadConfig reloads and updates the existing configuration
// There is only parameters that can be changed without restart.
// For example, if you want to change address or port,
// you need to restart the service anyway.
func (cfg *Config) ReloadConfig() error {
	newCfg, err := loadAndProcessConfig()
	if err != nil {
		return err
	}
	if cfg.GitHubToken != newCfg.GitHubToken {
		printConfigEvent("GitHub token has changed")
		cfg.GitHubToken = newCfg.GitHubToken
	}
	if cfg.LogLevel != newCfg.LogLevel {
		printConfigEventf("Log level has changed from %s to %s", cfg.LogLevel, newCfg.LogLevel)
		cfg.LogLevel = newCfg.LogLevel
	}
	if cfg.MemoryLimitBytes != newCfg.MemoryLimitBytes {
		printConfigEventf("Memory limit has changed from %d to %d (bytes)", cfg.MemoryLimitBytes, newCfg.MemoryLimitBytes)
		cfg.MemoryLimitBytes = newCfg.MemoryLimitBytes
	}
	if cfg.MemoryTTL != newCfg.MemoryTTL {
		printConfigEventf("Memory TTL has changed from %s to %s (seconds)", cfg.MemoryTTL, newCfg.MemoryTTL)
		cfg.MemoryTTL = newCfg.MemoryTTL
	}
	if cfg.PushMetricsUrl != newCfg.PushMetricsUrl {
		printConfigEventf("Push metrics address has changed to %s", newCfg.PushMetricsUrl)
		cfg.PushMetricsUrl = newCfg.PushMetricsUrl
	}
	if cfg.WebhookSecret != newCfg.WebhookSecret {
		printConfigEvent("Webhook secret has changed")
		cfg.WebhookSecret = newCfg.WebhookSecret
	}

	printConfigEvent("Configuration reloaded successfully")
	return nil
}

// parseSize parses a human-readable size string (e.g. "10G", "512M") into bytes
func parseSize(sizeStr string) (uint64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	if sizeStr == "" {
		return 0, fmt.Errorf("size string is empty")
	}

	units := map[string]uint64{
		"B":  1,
		"K":  1 << 10,
		"KB": 1 << 10,
		"M":  1 << 20,
		"MB": 1 << 20,
		"G":  1 << 30,
		"GB": 1 << 30,
		"T":  1 << 40,
		"TB": 1 << 40,
	}

	var numPart, unitPart string
	for i, r := range sizeStr {
		if (r < '0' || r > '9') && r != '.' {
			numPart = strings.TrimSpace(sizeStr[:i])
			unitPart = strings.TrimSpace(sizeStr[i:])
			break
		}
	}

	if numPart == "" {
		numPart = sizeStr
	}

	numValue, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size number: %v", err)
	}

	unitPart = strings.ToUpper(unitPart)
	if unitPart == "" {
		unitPart = "B"
	}

	multiplier, ok := units[unitPart]
	if !ok {
		return 0, fmt.Errorf("unrecognized size unit: %s", unitPart)
	}

	return uint64(numValue * float64(multiplier)), nil
}
