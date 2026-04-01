package tasks_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/equixss/todo-web/internal/core/domain"
)

func ph(idx int) string { return fmt.Sprintf("$%d", idx) }

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	limit *int,
	offset *int,
	userID *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	base := `
        SELECT id, version, title, description, completed,
               created_at, completed_at, user_id
        FROM todoapp.tasks
    `

	var where []string
	var args []interface{}

	if userID != nil {
		where = append(where, fmt.Sprintf("user_id = %s", ph(len(args)+1)))
		args = append(args, *userID)
	}
	limitPl := ph(len(args) + 1)
	offsetPl := ph(len(args) + 2)
	args = append(args, limit, offset)

	var sb strings.Builder
	sb.WriteString(strings.TrimSpace(base))
	if len(where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(where, " AND "))
	}
	sb.WriteString(fmt.Sprintf(" ORDER BY id ASC LIMIT %s OFFSET %s;", limitPl, offsetPl))

	query := sb.String()
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks error: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var model TaskModel
		err := rows.Scan(
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
			return nil, fmt.Errorf("scan tasks error: %w", err)
		}
		taskModels = append(taskModels, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows error: %w", err)
	}

	tasksDomain := taskDomainsFromModels(taskModels)
	return tasksDomain, nil
}
