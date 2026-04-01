package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, name, phone
	FROM todoapp.users
	ORDER BY id ASC 
	LIMIT $1
	OFFSET $2;
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users error: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel
	for rows.Next() {
		var model UserModel
		err := rows.Scan(
			&model.ID,
			&model.Version,
			&model.Name,
			&model.Phone,
		)
		if err != nil {
			return nil, fmt.Errorf("scan users error: %w", err)
		}
		userModels = append(userModels, model)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows error: %w", err)
	}

	usersDomain := userDomainsFromModels(userModels)
	return usersDomain, nil
}
