package handlers

import (
	"net/http"

	"github.com/Melsoft-Games/ant-watcher/internal/config"
	"github.com/Melsoft-Games/ant-watcher/internal/logger"
	"github.com/Melsoft-Games/ant-watcher/internal/store"
	"github.com/google/go-github/v66/github"
)

// WebhookHandler обрабатывает входящие запросы GitHub Webhook
type WebhookHandler struct {
	Store  *store.Store
	Secret []byte
}

// NewWebhookHandler инициализирует хендлер для вебхуков
func NewWebhookHandler(store *store.Store, cfg *config.Config) http.Handler {
	return &WebhookHandler{
		Store:  store,
		Secret: []byte(cfg.WebhookSecret),
	}
}

// ServeHTTP обрабатывает запросы вебхуков
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// payload, err := github.ValidatePayload(r, h.Secret)
	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		logger.Errorf("Invalid payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		logger.Errorf("Failed to parse webhook: %v", err)
		http.Error(w, "Failed to parse webhook", http.StatusInternalServerError)
		return
	}

	switch e := event.(type) {
	// case *github.CheckRunEvent:
	// 	h.handleCheckRun(e)
	// case *github.CheckSuiteEvent:
	// 	h.handleCheckSuite(e)
	case *github.WorkflowRunEvent:
		h.handleWorkflowRun(e)
	case *github.WorkflowJobEvent:
		h.handleWorkflowJob(e)
	case *github.WorkflowDispatchEvent:
		h.handleWorkflowDispatch(e)
	default:
		logger.Infof("Unhandled event type: %s", github.WebHookType(r))
		w.WriteHeader(http.StatusOK)
	}
}

// handleWorkflowRun обрабатывает событие WorkflowRunEvent
func (h *WebhookHandler) handleWorkflowRun(event *github.WorkflowRunEvent) {
	runRaw := event.WorkflowRun

	logger.Debug("WorkflowRunEvent: ", runRaw)

    // Проверка, что это действительно WorkflowRun
    if runRaw == nil {
        logger.Errorf("WorkflowRun is nil")
        return
    }

	runID := runRaw.GetID()
	runNumber := runRaw.GetRunNumber()
	status := runRaw.GetStatus()

 	// Проверяем статус на наличие
    if status == "" {
        logger.Errorf("Status is nil or empty for RunID=%d", runID)
        return
    }

	workflowRun := &github.WorkflowRun{
	    ID:           github.Int64(runID),             // Проверяем, что runID не nil
	    RunNumber:    &runNumber,                      // Убираем github.Int, если runNumber уже int
		Status:       github.String(status),
	    Conclusion:   runRaw.Conclusion,
	    CreatedAt:    runRaw.CreatedAt,
	    RunStartedAt: runRaw.RunStartedAt,
	    UpdatedAt:    runRaw.UpdatedAt,
	    Actor:        runRaw.Actor,
	}

	// Обновляем хранилище воркфлоу-ранов
	//h.Store.AddOrUpdateWorkflowRun(runID, workflowRun)
	h.Store.AddOrUpdateWorkflowRun(event.WorkflowRun.GetID(), workflowRun)
	h.Store.AddOrUpdateOrganization(event.Org.GetID(), event.Org)
	h.Store.AddOrUpdateRepository(event.Repo.GetID(), event.Repo)

	logger.Infof("WorkflowRun handled: RunID=%d, RunNumber=%d, Status=%s, Conclusion=%s, CreatedAt=%s, RunStartedAt=%s",
		runID, runNumber, runRaw.GetStatus(), runRaw.GetConclusion(), runRaw.GetCreatedAt(), runRaw.GetRunStartedAt())
}

// handleWorkflowJob обрабатывает событие WorkflowJobEvent
func (h *WebhookHandler) handleWorkflowJob(event *github.WorkflowJobEvent) {
	job := event.WorkflowJob
	// Ищем запуск воркфлоу
	runID := job.GetRunID()

	// Обработка времени создания, начала и завершения джоба
	jobModel := &github.WorkflowJob{
		ID:          github.Int64(job.GetID()),
		Status:      job.Status,
		Conclusion:  job.Conclusion,
		CreatedAt:   job.CreatedAt,  // Добавляем время создания
		StartedAt:   job.StartedAt,  // Добавляем время начала
		CompletedAt: job.CompletedAt,  // Добавляем время завершения
		RunnerID:    job.RunnerID,
		RunnerName:  job.RunnerName,
		Labels:      job.Labels,
		Steps:       job.Steps,
	}

	// Добавляем или обновляем джоб
	h.Store.AddOrUpdateJob(runID, jobModel)

	logger.Infof("WorkflowJob handled: JobID=%d, RunID=%d, Status=%s, Conclusion=%s, CreatedAt=%s, StartedAt=%s, CompletedAt=%s",
		job.GetID(), runID, job.GetStatus(), job.GetConclusion(), job.GetCreatedAt(), job.GetStartedAt(), job.GetCompletedAt())
}

// handleWorkflowDispatch обрабатывает событие WorkflowDispatchEvent
func (h *WebhookHandler) handleWorkflowDispatch(event *github.WorkflowDispatchEvent) {
	logger.Infof("Workflow dispatch event triggered")
}

// Вспомогательные функции для создания указателей и хеш-функции
func int64Ptr(i int64) *int64    { return &i }
func stringPtr(s string) *string { return &s }

// hashString простая хеш-функция для строк
func hashString(s string) int {
	var hash int
	for _, c := range s {
		hash += int(c)
	}
	return hash
}
