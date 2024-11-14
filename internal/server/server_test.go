package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Melsoft-Games/ant-watcher/internal/config"
	"github.com/Melsoft-Games/ant-watcher/internal/server"
	"github.com/Melsoft-Games/ant-watcher/internal/store"
	"github.com/stretchr/testify/assert"
)

// TestNewServer проверяет, что серверы инициализируются корректно
func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		WebhookAddress: "127.0.0.1",
		WebhookPort:    "8080",
		AdminAddress:   "127.0.0.1",
		AdminPort:      "8081",
		MetricsAddress: "127.0.0.1",
		MetricsPort:    "8082",
	}

	dataStore := store.NewStore()

	srv := server.NewServer(cfg, dataStore)

	assert.NotNil(t, srv)
	assert.NotNil(t, srv.WebhookMux)
	assert.NotNil(t, srv.AdminMux)
	assert.NotNil(t, srv.MetricsMux)
}

// // TestWebhookServer проверяет обработку запросов для вебхуков
// func TestWebhookServer(t *testing.T) {
// 	cfg := &config.Config{}
// 	dataStore := store.NewStore()

// 	srv := server.NewServer(cfg, dataStore)

// 	// Создаем тестовый HTTP-сервер
// 	testServer := httptest.NewServer(srv.WebhookMux)
// 	defer testServer.Close()

// 	// Отправляем тестовый запрос
// 	resp, err := http.Get(testServer.URL + "/webhook")
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)
// }

// TestAdminServer проверяет обработку запросов для администрирования
func TestAdminServer(t *testing.T) {
	cfg := &config.Config{}
	dataStore := store.NewStore()

	srv := server.NewServer(cfg, dataStore)

	// Создаем тестовый HTTP-сервер
	testServer := httptest.NewServer(srv.AdminMux)
	defer testServer.Close()

	// Отправляем тестовый запрос
	resp, err := http.Get(testServer.URL + "/status")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestMetricsServer проверяет обработку запросов для метрик
func TestMetricsServer(t *testing.T) {
	cfg := &config.Config{}
	dataStore := store.NewStore()

	srv := server.NewServer(cfg, dataStore)

	// Создаем тестовый HTTP-сервер
	testServer := httptest.NewServer(srv.MetricsMux)
	defer testServer.Close()

	// Отправляем тестовый запрос
	resp, err := http.Get(testServer.URL + "/metrics")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
