package server

import (
	"net/http"

	"github.com/Melsoft-Games/ant-watcher/internal/config"
	"github.com/Melsoft-Games/ant-watcher/internal/handlers"
	"github.com/Melsoft-Games/ant-watcher/internal/store"
)

// Server представляет HTTP серверы
type Server struct {
	Config     *config.Config
	Store      *store.Store
	WebhookMux *http.ServeMux
	AdminMux   *http.ServeMux
	MetricsMux *http.ServeMux
}

// NewServer initializes a new Server
func NewServer(cfg *config.Config, s *store.Store) *Server {
	// Мультиплексор для вебхуков
	webhookMux := http.NewServeMux()
	webhookHandler := handlers.NewWebhookHandler(s, cfg)
	webhookMux.Handle("/", webhookHandler)

	// Мультиплексор для административных маршрутов
	adminMux := http.NewServeMux()
	adminHandler := handlers.NewAdminHandler(cfg, s)
	adminMux.Handle("/", adminHandler)

	// Мультиплексор для метрик
	metricsMux := http.NewServeMux()
	metricsHandler := handlers.NewMetricsHandler(cfg)
	metricsMux.Handle("/metrics", metricsHandler)

	return &Server{
		Config:     cfg,
		Store:      s,
		WebhookMux: webhookMux,
		AdminMux:   adminMux,
		MetricsMux: metricsMux,
	}
}

// StartWebhookServer запускает сервер для вебхуков
func (srv *Server) StartWebhookServer(addr string) error {
	return http.ListenAndServe(addr, srv.WebhookMux)
}

// StartAdminServer запускает сервер для административных хендлеров
func (srv *Server) StartAdminServer(addr string) error {
	return http.ListenAndServe(addr, srv.AdminMux)
}

// StartMetricsServer запускает сервер для метрик
func (srv *Server) StartMetricsServer(addr string) error {
	return http.ListenAndServe(addr, srv.MetricsMux)
}

// Shutdown graceful shutdown серверов (если нужно)
func (srv *Server) Shutdown() error {
	// Логику graceful shutdown можно добавить, если нужно
	return nil
}
