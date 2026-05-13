package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_postgres_pool "github.com/equixss/todo-web/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.users
		SET name = $1, phone = $2, email = $3, version = version+1
		WHERE id = $4 AND version = $5
		RETURNING id, version, name, phone, email, COALESCE(password_hash, '') as password_hash;
	`

	row := r.pool.QueryRow(ctx, query, user.Name, user.Phone, user.Email, id, user.Version)

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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently accessed:%w",
				id,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
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
