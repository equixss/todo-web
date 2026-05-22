package tasks_service

import (
	"context"
	"fmt"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
)

func (s *TasksService) DeleteTask(ctx context.Context, id int, userID int) error {
	if id <= 0 {
		return fmt.Errorf("invalid task id: %w", core_errors.ErrInvalidArgument)
	}

	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	if task.AuthorUserID != userID {
		return fmt.Errorf("access denied: %w", core_errors.ErrNotFound)
	}

	if err := s.tasksRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
