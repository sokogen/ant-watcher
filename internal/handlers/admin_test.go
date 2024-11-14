package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Melsoft-Games/ant-watcher/internal/config"
	"github.com/Melsoft-Games/ant-watcher/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler(t *testing.T) {
	// Create mocks for configuration
	mockConfig := &config.Config{
		AdminAddress:       "127.0.0.1",
		AdminPort:          "8081",
		DisableAdminServer: "false",
		DisableAPI:         "false",
		FetchHistory:       "true",
		GitHubAPIURL:       "https://api-server.github.com",
		GitHubToken:        "example_token",
		LogLevel:           "DEBUG",
		MemoryLimit:        "512MB",
		MemoryTTL:          "24h",
		MetricsAddress:     "127.0.0.1",
		MetricsPort:        "9090",
		PushMetricsUrl:     "http://metrics:9091",
		WebhookAddress:     "127.0.0.1",
		WebhookPort:        "8082",
		WebhookSecret:      "supersecret",
	}

	// Create a mock store
	mockStore := &store.Store{}

	// Initialize the handler
	handler := NewAdminHandler(mockConfig, mockStore)

	// Check that the handler returns the correct response to the /admin/print-config request
	req := httptest.NewRequest(http.MethodGet, "/admin/print-config", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t,
		`{"admin_address":"127.0.0.1",
		"admin_port":"8081",
		"disable_admin_server":"false",
		"disable_api":"false",
		"fetch_history":"true",
		"github_api_url":"https://api-server.github.com",
		"github_token":"example_token",
		"log_level":"DEBUG",
		"memory_limit":"512MB",
		"memory_ttl":"24h",
		"metrics_address":"127.0.0.1",
		"metrics_port":"9090",
		"push_metrics_url":"http://metrics:9091",
		"webhook_address":"127.0.0.1",
		"webhook_port":"8082",
		"webhook_secret":"supersecret"}`,
		w.Body.String())

	// Check that the handler returns the correct response to the /status request
	req = httptest.NewRequest(http.MethodGet, "/status", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())

	// Checking /admin/reload-config
	req = httptest.NewRequest(http.MethodGet, "/admin/reload-config", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}
