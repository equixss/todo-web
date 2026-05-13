package tasks_service

import (
	"context"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	id int,
	userID int,
) (domain.Task, error) {
	if id <= 0 {
		return domain.Task{}, fmt.Errorf("invalid task id: %w", core_errors.ErrInvalidArgument)
	}

	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	if task.AuthorUserID != userID {
		return domain.Task{}, fmt.Errorf("access denied: %w", core_errors.ErrNotFound)
	}

	return task, nil
}
