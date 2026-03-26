package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO todoapp.users (name, phone)
		VALUES ($1,$2)
		RETURNING id, version, name, phone;
		`
	row := r.pool.QueryRow(ctx, query, user.Name, user.Phone)

	var model UserModel

	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.Phone,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return domain.NewUser(
		model.ID,
		model.Version,
		model.Name,
		model.Phone,
	), nil
}
