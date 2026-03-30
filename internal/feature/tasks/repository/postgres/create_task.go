package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_postgres_pool "github.com/equixss/todo-web/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
	INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, user_id)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id,version,title,description,completed,created_at,completed_at,user_id;
`
	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)
	var model TaskModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Title,
		&model.Description,
		&model.Completed,
		&model.CreatedAt,
		&model.CompletedAt,
		&model.AuthorUserID,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf(
				"%v: user sith id %d : %w",
				err,
				model.AuthorUserID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("create task scan error: %w", err)
	}

	return domain.NewTask(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.Completed,
		model.CreatedAt,
		model.CompletedAt,
		model.AuthorUserID,
	), nil
}
