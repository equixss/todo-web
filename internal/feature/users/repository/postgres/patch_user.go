package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	"github.com/jackc/pgx/v5"
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
		SET name = $1, phone = $2, version = version+1
		WHERE id = $3 AND version = $4
		RETURNING id, version, name, phone;
	`

	row := r.pool.QueryRow(ctx, query, user.Name, user.Phone, id, user.Version)

	var model UserModel
	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently accessed:%w",
				id,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}
	userDomain := domain.User{
		ID:      model.ID,
		Version: model.Version,
		Name:    model.Name,
		Phone:   model.Phone,
	}
	return userDomain, nil
}
