package tasks_service

import (
	"context"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("invalid task: %w", err)
	}
	newTask, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task error:%w", err)
	}

	return newTask, nil
}
