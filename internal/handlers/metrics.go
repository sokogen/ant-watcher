// internal/handlers/metrics.go
// metrics handler

package handlers

import (
	"net/http"

	"github.com/Melsoft-Games/ant-watcher/internal/config"
)

// NewMetricsHandler инициализирует хендлер для метрик Prometheus/VictoriaMetrics
func NewMetricsHandler(cfg *config.Config) http.Handler {
	// Здесь будет логика хендлера для метрик
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Metrics handler"))
	})
}
