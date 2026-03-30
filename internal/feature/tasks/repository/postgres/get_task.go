package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_postgres_pool "github.com/equixss/todo-web/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id,version,title,description,completed,created_at,completed_at,user_id
	FROM todoapp.tasks
	WHERE id = $1;
`
	row := r.pool.QueryRow(ctx, query, id)
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
				"task with id='%d' not found:%w",
				id,
				core_errors.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error:%w", err)
	}

	return domain.Task{
		ID:           model.ID,
		Version:      model.Version,
		Title:        model.Title,
		Description:  model.Description,
		Completed:    model.Completed,
		CreatedAt:    model.CreatedAt,
		CompletedAt:  model.CompletedAt,
		AuthorUserID: model.AuthorUserID,
	}, nil
}
