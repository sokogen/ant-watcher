package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Melsoft-Games/ant-watcher/internal/config"
	"github.com/Melsoft-Games/ant-watcher/internal/logger"
	"github.com/Melsoft-Games/ant-watcher/internal/store"
)

// AdminHandler отвечает за административные функции сервиса
type AdminHandler struct {
	Config *config.Config
	Store  *store.Store
}

// NewAdminHandler инициализирует хендлер для административных операций
func NewAdminHandler(cfg *config.Config, s *store.Store) http.Handler {
	return &AdminHandler{
		Config: cfg,
		Store:  s,
	}
}

// ServeHTTP обрабатывает запросы на административные функции
func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/status":
		h.handleStatus(w, r)
	case "/admin/reload-config":
		h.handleReloadConfig(w, r)
	case "/admin/print-config":
		h.handlePrintConfig(w, r)
	case "/admin/get-store":
		h.handleGetStore(w, r)
	case "/admin/organizations":
		h.handleGetAllOrganizations(w, r)
	case "/admin/repositories":
		h.handleGetAllRepositories(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// handleStatus возвращает статус сервиса
func (h *AdminHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleReloadConfig перезагружает конфигурацию
func (h *AdminHandler) handleReloadConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("Reloading initiated via admin endpoint")

	if err := h.Config.ReloadConfig(); err != nil {
		logger.Errorf("Failed to reload configuration: %v", err)
		http.Error(w, "Failed to reload configuration", http.StatusInternalServerError)
		return
	}

	logger.Info("Configuration reloaded successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configuration reloaded successfully"))
}

// handlePrintConfig выводит текущую конфигурацию в формате JSON
func (h *AdminHandler) handlePrintConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	configData, err := json.Marshal(h.Config)
	if err != nil {
		logger.Errorf("Failed to marshal configuration: %v", err)
		http.Error(w, "Failed to print configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(configData)
}

// handleGetOrganizations выводит все организации и их содержимое
func (h *AdminHandler) handleGetAllOrganizations(w http.ResponseWriter, r *http.Request) {
	h.Store.Mu.RLock()
	defer h.Store.Mu.RUnlock()

	organizations := h.Store.Organizations
	response, err := json.MarshalIndent(organizations, "", "  ")
	if err != nil {
		logger.Errorf("Failed to marshal organizations: %v", err)
		http.Error(w, "Failed to retrieve organizations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// handleGetRepositories выводит все репозитории в указанной организации
func (h *AdminHandler) handleGetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repos := h.Store.GetAllRepositories()
	if repos == nil {
		http.Error(w, "Failed to retrieve repositories", http.StatusInternalServerError)
		return
	}

	response, err := json.MarshalIndent(repos, "", "  ")
	if err != nil {
		logger.Errorf("Failed to marshal repositories: %v", err)
		http.Error(w, "Failed to retrieve repositories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// handleGetStore выводит информацию о хранилище
func (h *AdminHandler) handleGetStore(w http.ResponseWriter, r *http.Request) {
	h.Store.Mu.RLock()
	defer h.Store.Mu.RUnlock()

	logger.Infof("Users: %+v", h.Store.Users)
	logger.Infof("Organizations: %+v", h.Store.Organizations)

	response, err := json.MarshalIndent(h.Store, "", "  ")
	if err != nil {
		logger.Errorf("Failed to marshal store: %v", err)
		http.Error(w, "Failed to retrieve store", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
