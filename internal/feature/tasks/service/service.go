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
	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)
	DeleteTask(
		ctx context.Context,
		id int,
	) error
	GetTasks(
		ctx context.Context,
		limit *int,
		offset *int,
		userID *int,
	) ([]domain.Task, error)
}

type TasksService struct {
	tasksRepository TasksRepository
}

func NewTasksService(repository TasksRepository) *TasksService {
	return &TasksService{
		tasksRepository: repository,
	}
}
