package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Melsoft-Games/ant-watcher/internal/config"
	"github.com/Melsoft-Games/ant-watcher/internal/logger"
	"github.com/Melsoft-Games/ant-watcher/internal/server"
	"github.com/Melsoft-Games/ant-watcher/internal/store"
)

func main() {
	// Initialize the logger
	logger.Init()

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set the log level
	logger.ChangeLogLevel(cfg.LogLevel)

	// Initialize the data store
	dataStore := store.NewStore()

	// Create a new server
	srv := server.NewServer(cfg, dataStore)

	// Запуск админ-сервера
	go func() {
		adminAddr := fmt.Sprintf("%s:%s", cfg.AdminAddress, cfg.AdminPort)
		logger.Infof("Starting admin server on %s", adminAddr)
		if err := srv.StartAdminServer(adminAddr); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Admin server failed: %v", err)
		}
	}()

	// Запуск сервера для вебхуков
	go func() {
		webhookAddr := fmt.Sprintf("%s:%s", cfg.WebhookAddress, cfg.WebhookPort)
		logger.Infof("Starting webhook server on %s", webhookAddr)
		if err := srv.StartWebhookServer(webhookAddr); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Webhook server failed: %v", err)
		}
	}()

	// Запуск сервера для метрик (если требуется)
	go func() {
		metricsAddr := fmt.Sprintf("%s:%s", cfg.MetricsAddress, cfg.MetricsPort)
		logger.Infof("Starting metrics server on %s", metricsAddr)
		if err := srv.StartMetricsServer(metricsAddr); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Metrics server failed: %v", err)
		}
	}()

	// Wait for the shutdown signal, to gracefully shutdown the servers
	waitForShutdown(srv)
}

// gracefully shutdown the servers
func waitForShutdown(srv *server.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	// timeout for graceful shutdown by default is 5 seconds
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call graceful shutdown for the servers
	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
