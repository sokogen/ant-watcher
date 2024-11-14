package api

import (
	"context"
	"time"

	"github.com/Melsoft-Games/ant-watcher/internal/store"
	"github.com/google/go-github/v66/github"
)

// FetchRecentWorkflows собирает данные о WorkflowRun за указанный период из всех репозиториев организации
func FetchRecentWorkflows(client *github.Client, org string, store *store.Store, fetchHistory time.Duration) error {
	// Получаем текущую дату и вычисляем дату начала
	fromTime := time.Now().Add(-fetchHistory)

	// Получаем список всех репозиториев организации
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opt)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			err := syncRepoWorkflows(client, org, *repo.Name, store, fromTime)
			if err != nil {
				return err
			}
		}

		// Если есть следующая страница, продолжаем
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return nil
}

// syncRepoWorkflows синхронизирует WorkflowRun репозитория за определённый период
func syncRepoWorkflows(client *github.Client, owner, repo string, store *store.Store, fromTime time.Time) error {
	opt := &github.ListWorkflowRunsOptions{
		Created: fromTime.Format(time.RFC3339), // Форматируем время в строку ISO 8601
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	}

	for {
		runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(context.Background(), owner, repo, opt)
		if err != nil {
			return err
		}
		// Добавляем или обновляем каждый WorkflowRun в store
		for _, run := range runs.WorkflowRuns {
			store.AddOrUpdateWorkflow(run.GetID(), run)
		}
		// Переход на следующую страницу
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil
}
