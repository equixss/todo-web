package statistics_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/equixss/todo-web/internal/core/domain"
)

func ph(idx int) string { return fmt.Sprintf("$%d", idx) }

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	base := `
	SELECT id,version,title,description,completed,created_at,completed_at,user_id
	FROM todoapp.tasks
`
	var where []string
	var args []interface{}
	if userID != nil {
		where = append(where, fmt.Sprintf("user_id = %s", ph(len(args)+1)))
		args = append(args, *userID)
	}
	if from != nil {
		where = append(where, fmt.Sprintf("created_at >= %s", ph(len(args)+1)))
		args = append(args, *from)
	}
	if to != nil {
		where = append(where, fmt.Sprintf("created_at < %s", ph(len(args)+1)))
		args = append(args, *to)
	}
	var sb strings.Builder
	sb.WriteString(strings.TrimSpace(base))
	if len(where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(where, " AND "))
	}
	query := sb.String()
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks error: %w", err)
	}
	defer rows.Close()

	var tasksModels []TaskModel
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
			return nil, fmt.Errorf("scan tasks row error: %w", err)
		}
		tasksModels = append(tasksModels, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	tasksDomain := taskDomainsFromModels(tasksModels)
	return tasksDomain, nil
}
