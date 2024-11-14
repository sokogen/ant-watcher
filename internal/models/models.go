package models

import (
	"time"
)

// Step представляет этап выполнения внутри джоба
type Step struct {
	ID          *int64     `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Conclusion  *string    `json:"conclusion,omitempty"`
	Number      *int64     `json:"number,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Job представляет задачу в рамках запуска воркфлоу
type Job struct {
	ID              *int64          `json:"id,omitempty"`
	Name            *string         `json:"name,omitempty"`
	Status          *string         `json:"status,omitempty"`
	Conclusion      *string         `json:"conclusion,omitempty"`
	StartedAt       *time.Time      `json:"started_at,omitempty"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
	RunnerID        *int64          `json:"runner_id,omitempty"`
	RunnerName      *string         `json:"runner_name,omitempty"`
	RunnerOS        *string         `json:"runner_os,omitempty"`
	RunnerGroupID   *int64          `json:"runner_group_id,omitempty"`
	RunnerGroupName *string         `json:"runner_group_name,omitempty"`
	Steps           map[int64]*Step `json:"steps,omitempty"` // Этапы внутри джоба
}

// WorkflowRun представляет запуск воркфлоу с возможными перезапусками
type WorkflowRun struct {
	RunID       *int64         `json:"run_id,omitempty"`
	WorkflowID  *int64         `json:"workflow_id,omitempty"`
	RunNumber   *int64         `json:"run_number,omitempty"`
	Attempt     *int64         `json:"attempt,omitempty"` // Счетчик попыток запуска (перезапусков)
	Event       *string        `json:"event,omitempty"`
	Status      *string        `json:"status,omitempty"`
	Conclusion  *string        `json:"conclusion,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	TriggeredBy *User          `json:"triggered_by,omitempty"` // Инициатор запуска или перезапуска
	Jobs        map[int64]*Job `json:"jobs,omitempty"`         // Джобы внутри запуска
}

type Workflow struct {
	ID        *int64                 `json:"id,omitempty"`
	Name      *string                `json:"name,omitempty"`
	Path      *string                `json:"path,omitempty"`
	State     *string                `json:"state,omitempty"` // "active" или "disabled"
	CreatedAt *time.Time             `json:"created_at,omitempty"`
	UpdatedAt *time.Time             `json:"updated_at,omitempty"`
	URL       *string                `json:"url,omitempty"`
	HTMLURL   *string                `json:"html_url,omitempty"`
	BadgeURL  *string                `json:"badge_url,omitempty"`
	Runs      map[int64]*WorkflowRun `json:"runs,omitempty"` // Запуски воркфлоу
}

type Repository struct {
	ID          *int64              `json:"id,omitempty"`
	Name        *string             `json:"name,omitempty"`
	FullName    *string             `json:"full_name,omitempty"`
	Owner       Owner               `json:"owner,omitempty"`
	Private     *bool               `json:"private,omitempty"`
	HTMLURL     *string             `json:"html_url,omitempty"`
	Description *string             `json:"description,omitempty"`
	Fork        *bool               `json:"fork,omitempty"`
	URL         *string             `json:"url,omitempty"`
	Workflows   map[int64]*Workflow `json:"workflows,omitempty"` // Воркфлоу, принадлежащие репозиторию
	Commits     map[string]*Commit  `json:"commits,omitempty"`   // Коммиты, связанные с репозиторием
	Releases    map[string]*Release `json:"releases,omitempty"`  // Релизы, связанные с репозиторием
}

// Commit представляет коммит
type Commit struct {
	SHA     *string    `json:"sha,omitempty"`
	Message *string    `json:"message,omitempty"`
	Author  *string    `json:"author,omitempty"`
	Date    *time.Time `json:"date,omitempty"`
}

// Release представляет релиз или тег
type Release struct {
	TagName    *string    `json:"tag_name,omitempty"`
	ReleasedAt *time.Time `json:"released_at,omitempty"`
}

// User представляет пользователя GitHub
type User struct {
	ID           *int64                `json:"id,omitempty"`
	Login        *string               `json:"login,omitempty"`
	Name         *string               `json:"name,omitempty"`
	Email        *string               `json:"email,omitempty"`
	Repositories map[int64]*Repository `json:"repositories,omitempty"` // Репозитории, принадлежащие пользователю
}

// Organization представляет организацию
type Organization struct {
	ID           *int64                `json:"id,omitempty"`
	Login        *string               `json:"login,omitempty"`
	Name         *string               `json:"name,omitempty"`
	URL          *string               `json:"url,omitempty"`
	HTMLURL      *string               `json:"html_url,omitempty"`
	Description  *string               `json:"description,omitempty"`
	Members      map[int64]*User       `json:"members,omitempty"`      // Члены организации
	Repositories map[int64]*Repository `json:"repositories,omitempty"` // Репозитории, принадлежащие организации
}

// Owner интерфейс для пользователей и организаций, представляющий владельца репозитория
type Owner interface {
	GetID() *int64
	GetLogin() *string
	GetRepositories() map[int64]*Repository
}

// Реализация интерфейса Owner для организации
func (o *Organization) GetID() *int64 {
	return o.ID
}

func (o *Organization) GetLogin() *string {
	return o.Login
}

func (o *Organization) GetRepositories() map[int64]*Repository {
	return o.Repositories
}

// Реализация интерфейса Owner для пользователя
func (u *User) GetID() *int64 {
	return u.ID
}

func (u *User) GetLogin() *string {
	return u.Login
}

func (u *User) GetRepositories() map[int64]*Repository {
	return u.Repositories
}
