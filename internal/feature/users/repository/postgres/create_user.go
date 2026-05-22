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
		INSERT INTO todoapp.users (name, phone, email, password_hash)
		VALUES ($1,$2,$3,$4)
		RETURNING id, version, name, phone, email, password_hash;
		`
	row := r.pool.QueryRow(ctx, query, user.Name, user.Phone, user.Email, user.PasswordHash)

	var model UserModel

	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.Phone,
		&model.Email,
		&model.PasswordHash,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return domain.NewUser(
		model.ID,
		model.Version,
		model.Name,
		model.Phone,
		model.Email,
		model.PasswordHash,
	), nil
}
