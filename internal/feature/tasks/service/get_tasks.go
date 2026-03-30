package tasks_service

import (
	"context"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	limit *int,
	offset *int,
	userID *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be greater than or equal to 0:%w",
			core_errors.ErrInvalidArgument,
		)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be greater than or equal to 0:%w",
			core_errors.ErrInvalidArgument,
		)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, limit, offset, userID)
	if err != nil {
		return []domain.Task{}, fmt.Errorf("get tasks error:%w", err)
	}
	return tasks, nil
}
