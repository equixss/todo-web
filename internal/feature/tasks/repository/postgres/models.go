package tasks_postgres_repository

import (
	"time"

	"github.com/equixss/todo-web/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomainFromModel(model TaskModel) domain.Task {
	return domain.NewTask(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.Completed,
		model.CreatedAt,
		model.CompletedAt,
		model.AuthorUserID,
	)
}

func taskDomainsFromModels(models []TaskModel) []domain.Task {
	tasksDomain := make([]domain.Task, len(models))
	for i, model := range models {
		tasksDomain[i] = taskDomainFromModel(model)
	}
	return tasksDomain
}
