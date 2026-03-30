package tasks_service

import (
	"context"

	"github.com/equixss/todo-web/internal/core/domain"
)

type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
}

type TasksService struct {
	tasksRepository TasksRepository
}

func NewTasksService(repository TasksRepository) *TasksService {
	return &TasksService{
		tasksRepository: repository,
	}
}
