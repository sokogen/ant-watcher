// internal/store/store.go

package store

import (
	"sync"

	"github.com/Melsoft-Games/ant-watcher/internal/logger"
	"github.com/google/go-github/v66/github"
)

type Store struct {
	Mu            sync.RWMutex                   `json:"-"`
	Users         map[int64]*github.User         `json:"users"`
	Organizations map[int64]*github.Organization `json:"organizations"`
	Repositories  map[int64]*github.Repository   `json:"repositories"`
	Workflows     map[int64]*github.Workflow     `json:"workflows"`
	WorkflowRuns  map[int64]*github.WorkflowRun  `json:"workflow_runs"`
	Jobs          map[int64]*github.WorkflowJob  `json:"jobs"`
}

// NewStore инициализирует хранилище
func NewStore() *Store {
	return &Store{
		Users:         make(map[int64]*github.User),
		Organizations: make(map[int64]*github.Organization),
		Repositories:  make(map[int64]*github.Repository),
		Workflows:     make(map[int64]*github.Workflow),
		WorkflowRuns:  make(map[int64]*github.WorkflowRun),
		Jobs:          make(map[int64]*github.WorkflowJob),
	}
}

// AddOrUpdateUser добавляет или обновляет пользователя
func (s *Store) AddOrUpdateUser(userID int64, user *github.User) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Users[userID] = user
	logger.Infof("User with ID: %d added/updated", userID)
}

// AddOrUpdateOrganization добавляет или обновляет организацию
func (s *Store) AddOrUpdateOrganization(orgID int64, org *github.Organization) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Organizations[orgID] = org
	logger.Infof("Organization with ID: %d added/updated", orgID)
}

// AddOrUpdateRepository добавляет или обновляет репозиторий
func (s *Store) AddOrUpdateRepository(repoID int64, repo *github.Repository) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Repositories[repoID] = repo
	logger.Infof("Repository with ID: %d added/updated", repoID)
}

// AddOrUpdateWorkflow добавляет или обновляет воркфлоу
func (s *Store) AddOrUpdateWorkflow(workflowID int64, workflow *github.Workflow) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Workflows[workflowID] = workflow
	logger.Infof("Workflow with ID: %d added/updated", workflowID)
}

// AddOrUpdateWorkflowRun добавляет или обновляет запуск воркфлоу
func (s *Store) AddOrUpdateWorkflowRun(runID int64, run *github.WorkflowRun) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.WorkflowRuns[runID] = run
	logger.Infof("WorkflowRun with ID: %d added/updated", runID)
}

// AddOrUpdateJob добавляет или обновляет джоб
func (s *Store) AddOrUpdateJob(jobID int64, job *github.WorkflowJob) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Jobs[jobID] = job
	logger.Infof("Job with ID: %d added/updated", jobID)
}

// GetUser возвращает пользователя по его ID
func (s *Store) GetUser(userID int64) (*github.User, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	user, exists := s.Users[userID]
	return user, exists
}

// GetOrganization возвращает организацию по её ID
func (s *Store) GetOrganization(orgID int64) (*github.Organization, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	org, exists := s.Organizations[orgID]
	return org, exists
}

// GetRepository возвращает репозиторий по его ID
func (s *Store) GetRepository(repoID int64) (*github.Repository, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	repo, exists := s.Repositories[repoID]
	return repo, exists
}

// GetWorkflow возвращает воркфлоу по его ID
func (s *Store) GetWorkflow(workflowID int64) (*github.Workflow, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	workflow, exists := s.Workflows[workflowID]
	return workflow, exists
}

// GetWorkflowRun возвращает запуск воркфлоу по его ID
func (s *Store) GetWorkflowRun(runID int64) (*github.WorkflowRun, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	run, exists := s.WorkflowRuns[runID]
	return run, exists
}

// GetJob возвращает джоб по его ID
func (s *Store) GetJob(jobID int64) (*github.WorkflowJob, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	job, exists := s.Jobs[jobID]
	return job, exists
}

// GetAllRepositories возвращает все репозитории
func (s *Store) GetAllRepositories() map[int64]*github.Repository {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.Repositories
}

// GetAllWorkflows возвращает все воркфлоу
func (s *Store) GetAllWorkflows() map[int64]*github.Workflow {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.Workflows
}

// GetAllWorkflowRuns возвращает все запуски воркфлоу
func (s *Store) GetAllWorkflowRuns() map[int64]*github.WorkflowRun {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.WorkflowRuns
}

// GetAllJobs возвращает все джобы
func (s *Store) GetAllJobs() map[int64]*github.WorkflowJob {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.Jobs
}

// GetAllUsers возвращает всех пользователей
func (s *Store) GetAllUsers() map[int64]*github.User {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.Users
}

// GetAllOrganizations возвращает все организации
func (s *Store) GetAllOrganizations() map[int64]*github.Organization {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.Organizations
}
