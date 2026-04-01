package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_postgres_pool "github.com/equixss/todo-web/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	id int,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.tasks
	SET 
	    title = $1, 
	    description = $2, 
	    version = version+1, 
	    completed = $3,
	    completed_at = $4
	WHERE id = $5 and version = $6
	RETURNING 
	    id,
	    version,
	    title,
	    description,
	    completed,
	    created_at,
	    completed_at,
	    user_id;
`
	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		id,
		task.Version,
	)

	var model TaskModel
	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Title,
		&model.Description,
		&model.Completed,
		&model.CreatedAt,
		&model.CompletedAt,
		&model.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d' concurrently accessed:%w",
				id,
				core_errors.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}
	return taskDomainFromModel(model), nil
}
